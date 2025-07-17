package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/spf13/cobra"
	"github.com/IshaanNene/EliteCode-brew/internal/docker"
	"github.com/IshaanNene/EliteCode-brew/internal/utils"
)

type metadata struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Language    string `json:"language"`
	TimeLimit   int    `json:"time_limit"`
	MemoryLimit int    `json:"memory_limit"`
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run your solution",
	Long: `Run your solution in a sandboxed Docker environment.
This will:
1. Build a Docker image for your code
2. Run the code with sample test cases
3. Show execution time and memory usage`,
	RunE: func(cmd *cobra.Command, args []string) error {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting working directory: %v", err)
		}

		metadataBytes, err := os.ReadFile(filepath.Join(wd, "metadata.json"))
		if err != nil {
			return fmt.Errorf("error reading metadata file: %v", err)
		}

		var meta metadata
		if err := json.Unmarshal(metadataBytes, &meta); err != nil {
			return fmt.Errorf("error parsing metadata: %v", err)
		}

		dockerClient, err := docker.NewClient()
		if err != nil {
			return fmt.Errorf("error creating Docker client: %v", err)
		}
		defer dockerClient.Close()

		dockerfile, err := os.ReadFile(filepath.Join(wd, "Dockerfile"))
		if err != nil {
			return fmt.Errorf("error reading Dockerfile: %v", err)
		}

		sourceCode, err := os.ReadFile(filepath.Join(wd, "main."+utils.GetFileExtension(meta.Language)))
		if err != nil {
			return fmt.Errorf("error reading source code: %v", err)
		}

		files := map[string][]byte{
			"Dockerfile": dockerfile,
			"main." + utils.GetFileExtension(meta.Language): sourceCode,
		}

		buildCtx, err := docker.CreateBuildContext(files)
		if err != nil {
			return fmt.Errorf("error creating build context: %v", err)
		}

		fmt.Println("Building Docker image...")
		buildOptions := types.ImageBuildOptions{
			Dockerfile: "Dockerfile",
			Tags:       []string{fmt.Sprintf("elitecode/%s:%s", meta.ID, meta.Language)},
			Remove:     true,
		}

		if err := dockerClient.BuildImage(cmd.Context(), buildCtx, buildOptions); err != nil {
			return fmt.Errorf("error building Docker image: %v", err)
		}

		containerConfig := &container.Config{
			Image: fmt.Sprintf("elitecode/%s:%s", meta.ID, meta.Language),
			Cmd:   utils.GetRunCommand(meta.Language),
		}

		hostConfig := &container.HostConfig{
			Resources: container.Resources{
				Memory:    int64(meta.MemoryLimit) * 1024 * 1024, // Convert MB to bytes
				CPUPeriod: 100000,
				CPUQuota:  50000, // Limit to 0.5 CPU
			},
		}

		containerName := fmt.Sprintf("elitecode_%s_%s_%d", meta.ID, meta.Language, time.Now().Unix())

		fmt.Println("Running your solution...")
		output, executionTime, memoryUsage, err := dockerClient.RunContainer(
			cmd.Context(),
			containerConfig,
			hostConfig,
			containerName,
		)
		if err != nil {
			return fmt.Errorf("error running container: %v", err)
		}

		fmt.Println("\nOutput:")
		fmt.Println("----------------------------------------")
		fmt.Printf("%s", output)
		fmt.Println("----------------------------------------")
		fmt.Printf("\nExecution Time: %v\n", executionTime)
		fmt.Printf("Memory Usage: %.2f MB\n", memoryUsage)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
