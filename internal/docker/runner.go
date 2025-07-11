package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"elitecode/internal/storage"
	"elitecode/internal/utils"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

type Runner struct {
	client *client.Client
	config storage.DockerConfig
}

type RunOptions struct {
	Language    string
	Code        string
	Input       string
	TimeLimit   time.Duration
	MemoryLimit string
	CPULimit    string
	WorkDir     string
}

type RunResult struct {
	Output     string
	Error      string
	ExitCode   int
	Time       time.Duration
	Memory     int64
	Success    bool
	TimedOut   bool
	StartTime  time.Time
	EndTime    time.Time
}

var supportedLanguages = map[string]string{
	"c":          "elitecode/c:latest",
	"cpp":        "elitecode/cpp:latest",
	"python":     "elitecode/python:latest",
	"java":       "elitecode/java:latest",
	"javascript": "elitecode/javascript:latest",
	"go":         "elitecode/go:latest",
	"rust":       "elitecode/rust:latest",
}

func NewRunner() (*Runner, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	config := storage.GetConfig()
	return &Runner{
		client: cli,
		config: config.Docker,
	}, nil
}

func (r *Runner) RunCode(ctx context.Context, options RunOptions) (*RunResult, error) {
	result := &RunResult{
		StartTime: time.Now(),
	}

	// Check if language is supported
	imageName, exists := supportedLanguages[options.Language]
	if !exists {
		return nil, fmt.Errorf("unsupported language: %s", options.Language)
	}

	// Ensure image exists
	if err := r.ensureImage(ctx, imageName); err != nil {
		return nil, fmt.Errorf("failed to ensure image: %w", err)
	}

	// Create temporary directory for code execution
	tempDir, err := os.MkdirTemp("", "elitecode-run-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Write code to file
	codeFile := filepath.Join(tempDir, getCodeFileName(options.Language))
	if err := os.WriteFile(codeFile, []byte(options.Code), 0644); err != nil {
		return nil, fmt.Errorf("failed to write code file: %w", err)
	}

	// Write input to file
	inputFile := filepath.Join(tempDir, "input.txt")
	if err := os.WriteFile(inputFile, []byte(options.Input), 0644); err != nil {
		return nil, fmt.Errorf("failed to write input file: %w", err)
	}

	// Create and run container
	containerResp, err := r.createContainer(ctx, imageName, tempDir, options)
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}
	defer r.cleanupContainer(ctx, containerResp.ID)

	// Start container
	if err := r.client.ContainerStart(ctx, containerResp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	// Wait for container with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, options.TimeLimit)
	defer cancel()

	statusCh, errCh := r.client.ContainerWait(timeoutCtx, containerResp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return nil, fmt.Errorf("error waiting for container: %w", err)
		}
	case status := <-statusCh:
		result.ExitCode = int(status.StatusCode)
	case <-timeoutCtx.Done():
		result.TimedOut = true
		// Force kill the container
		r.client.ContainerKill(ctx, containerResp.ID, "SIGKILL")
	}

	result.EndTime = time.Now()
	result.Time = result.EndTime.Sub(result.StartTime)

	// Get container logs
	output, err := r.getContainerLogs(ctx, containerResp.ID)
	if err != nil {
		utils.Logger.Error("Failed to get container logs", "error", err)
	}

	// Split output into stdout and stderr
	r.parseOutput(output, result)

	// Get container stats
	if err := r.getContainerStats(ctx, containerResp.ID, result); err != nil {
		utils.Logger.Error("Failed to get container stats", "error", err)
	}

	result.Success = result.ExitCode == 0 && !result.TimedOut

	return result, nil
}

func (r *Runner) createContainer(ctx context.Context, imageName, tempDir string, options RunOptions) (container.ContainerCreateCreatedBody, error) {
	// Set default limits if not provided
	memoryLimit := options.MemoryLimit
	if memoryLimit == "" {
		memoryLimit = r.config.MemoryLimit
	}

	cpuLimit := options.CPULimit
	if cpuLimit == "" {
		cpuLimit = r.config.CPULimit
	}

	// Convert memory limit to bytes
	memoryBytes, err := parseMemoryLimit(memoryLimit)
	if err != nil {
		return container.ContainerCreateCreatedBody{}, fmt.Errorf("invalid memory limit: %w", err)
	}

	// Convert CPU limit to nano CPUs
	cpuNanos, err := parseCPULimit(cpuLimit)
	if err != nil {
		return container.ContainerCreateCreatedBody{}, fmt.Errorf("invalid CPU limit: %w", err)
	}

	containerConfig := &container.Config{
		Image:           imageName,
		WorkingDir:      "/workspace",
		Cmd:             getRunCommand(options.Language),
		AttachStdout:    true,
		AttachStderr:    true,
		NetworkDisabled: true, // Disable network for security
	}

	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: tempDir,
				Target: "/workspace",
			},
		},
		Resources: container.Resources{
			Memory:   memoryBytes,
			NanoCPUs: cpuNanos,
		},
		AutoRemove: true,
	}

	return r.client.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
}

func (r *Runner) ensureImage(ctx context.Context, imageName string) error {
	// Check if image exists locally
	_, _, err := r.client.ImageInspectWithRaw(ctx, imageName)
	if err == nil {
		return nil // Image exists
	}

	// Pull image if not found
	utils.Logger.Info("Pulling Docker image", "image", imageName)
	reader, err := r.client.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image %s: %w", imageName, err)
	}
	defer reader.Close()

	// Read the pull response (required for the pull to complete)
	_, err = io.Copy(io.Discard, reader)
	return err
}

func (r *Runner) getContainerLogs(ctx context.Context, containerID string) (string, error) {
	logs, err := r.client.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return "", err
	}
	defer logs.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(logs)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (r *Runner) parseOutput(output string, result *RunResult) {
	lines := strings.Split(output, "\n")
	var stdout, stderr []string

	for _, line := range lines {
		if strings.HasPrefix(line, "STDERR:") {
			stderr = append(stderr, strings.TrimPrefix(line, "STDERR:"))
		} else if line != "" {
			stdout = append(stdout, line)
		}
	}

	result.Output = strings.Join(stdout, "\n")
	result.Error = strings.Join(stderr, "\n")
}

func (r *Runner) getContainerStats(ctx context.Context, containerID string, result *RunResult) error {
	stats, err := r.client.ContainerStats(ctx, containerID, false)
	if err != nil {
		return err
	}
	defer stats.Body.Close()

	var statsJSON types.StatsJSON
	if err := json.NewDecoder(stats.Body).Decode(&statsJSON); err != nil {
		return err
	}

	result.Memory = int64(statsJSON.MemoryStats.Usage)
	return nil
}

func (r *Runner) cleanupContainer(ctx context.Context, containerID string) {
	r.client.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
		Force: true,
	})
}

func getCodeFileName(language string) string {
	switch language {
	case "c":
		return "main.c"
	case "cpp":
		return "main.cpp"
	case "python":
		return "main.py"
	case "java":
		return "Main.java"
	case "javascript":
		return "main.js"
	case "go":
		return "main.go"
	case "rust":
		return "main.rs"
	default:
		return "main.txt"
	}
}

func getRunCommand(language string) []string {
	switch language {
	case "c":
		return []string{"sh", "-c", "gcc -o main main.c && ./main < input.txt"}
	case "cpp":
		return []string{"sh", "-c", "g++ -o main main.cpp && ./main < input.txt"}
	case "python":
		return []string{"sh", "-c", "python3 main.py < input.txt"}
	case "java":
		return []string{"sh", "-c", "javac Main.java && java Main < input.txt"}
	case "javascript":
		return []string{"sh", "-c", "node main.js < input.txt"}
	case "go":
		return []string{"sh", "-c", "go run main.go < input.txt"}
	case "rust":
		return []string{"sh", "-c", "rustc main.rs && ./main < input.txt"}
	default:
		return []string{"echo", "Unsupported language"}
	}
}

func parseMemoryLimit(limit string) (int64, error) {
	if limit == "" {
		return 128 * 1024 * 1024, nil // 128MB default
	}

	multiplier := int64(1)
	unit := limit[len(limit)-1:]
	value := limit[:len(limit)-1]

	switch strings.ToLower(unit) {
	case "k":
		multiplier = 1024
	case "m":
		multiplier = 1024 * 1024
	case "g":
		multiplier = 1024 * 1024 * 1024
	default:
		value = limit
	}

	val, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}

	return val * multiplier, nil
}

func parseCPULimit(limit string) (int64, error) {
	if limit == "" {
		return 1000000000, nil // 1 CPU default
	}

	val, err := strconv.ParseFloat(limit, 64)
	if err != nil {
		return 0, err
	}

	return int64(val * 1000000000), nil // Convert to nano CPUs
}