package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Client struct {
	*client.Client
}

func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("error creating Docker client: %v", err)
	}
	return &Client{cli}, nil
}

func (c *Client) RunContainer(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, name string) (output []byte, executionTime time.Duration, memoryUsage float64, err error) {
	createResp, err := c.ContainerCreate(ctx, config, hostConfig, nil, nil, name)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("error creating container: %v", err)
	}

	startTime := time.Now()
	if err := c.ContainerStart(ctx, createResp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, 0, 0, fmt.Errorf("error starting container: %v", err)
	}

	statusCh, errCh := c.ContainerWait(ctx, createResp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return nil, 0, 0, fmt.Errorf("error waiting for container: %v", err)
		}
	case <-statusCh:
	}

	logs, err := c.ContainerLogs(ctx, createResp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return nil, 0, 0, fmt.Errorf("error getting container logs: %v", err)
	}
	defer logs.Close()

	output, err = io.ReadAll(logs)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("error reading logs: %v", err)
	}

	stats, err := c.ContainerStats(ctx, createResp.ID, false)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("error getting container stats: %v", err)
	}
	defer stats.Body.Close()

	var statsJSON types.StatsJSON
	if err := json.NewDecoder(stats.Body).Decode(&statsJSON); err != nil {
		return nil, 0, 0, fmt.Errorf("error decoding stats: %v", err)
	}

	executionTime = time.Since(startTime)
	memoryUsage = float64(statsJSON.MemoryStats.Usage) / 1024 / 1024 // Convert to MB

	if err := c.ContainerRemove(ctx, createResp.ID, types.ContainerRemoveOptions{
		Force: true,
	}); err != nil {
		return nil, 0, 0, fmt.Errorf("error removing container: %v", err)
	}

	return output, executionTime, memoryUsage, nil
}

func (c *Client) BuildImage(ctx context.Context, buildCtx io.Reader, options types.ImageBuildOptions) error {
	buildResponse, err := c.ImageBuild(ctx, buildCtx, options)
	if err != nil {
		return fmt.Errorf("error building Docker image: %v", err)
	}
	defer buildResponse.Body.Close()

	if _, err := io.Copy(os.Stdout, buildResponse.Body); err != nil {
		return fmt.Errorf("error reading build output: %v", err)
	}

	return nil
}

func CreateBuildContext(files map[string][]byte) (io.Reader, error) {
	buildCtx := new(bytes.Buffer)
	tw := tar.NewWriter(buildCtx)

	for name, content := range files {
		err := tw.WriteHeader(&tar.Header{
			Name: name,
			Size: int64(len(content)),
			Mode: 0644,
		})
		if err != nil {
			return nil, fmt.Errorf("error writing header for %s: %v", name, err)
		}

		if _, err := tw.Write(content); err != nil {
			return nil, fmt.Errorf("error writing content for %s: %v", name, err)
		}
	}

	if err := tw.Close(); err != nil {
		return nil, fmt.Errorf("error closing tar writer: %v", err)
	}

	return buildCtx, nil
}
