package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func RunInDocker(language, mountDir, testcases string, isSubmit bool) error {
	var dockerfile string
	switch language {
	case "cpp":
		dockerfile = "cpp.Dockerfile"
	case "java":
		dockerfile = "java.Dockerfile"
	case "python":
		dockerfile = "python.Dockerfile"
	case "go":
		dockerfile = "go.Dockerfile"
	default:
		return fmt.Errorf("unsupported language: %s", language)
	}

	containerName := "elitecode_run_container"
	tag := "elitecode_" + language + "_image"
	dockerfilePath := filepath.Join("docker", dockerfile)

	fmt.Println("🛠️ Building Docker image...")
	cmdBuild := exec.Command("docker", "build", "-t", tag, "-f", dockerfilePath, ".")
	cmdBuild.Stdout = os.Stdout
	cmdBuild.Stderr = os.Stderr
	if err := cmdBuild.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	fmt.Println("🐳 Running inside Docker...")
	cmdRun := exec.Command("docker", "run", "--rm", "-v", fmt.Sprintf("%s:/app", mountDir), "-w", "/app", tag)
	cmdRun.Stdout = os.Stdout
	cmdRun.Stderr = os.Stderr
	cmdRun.Stdin = os.Stdin

	if err := cmdRun.Run(); err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	if isSubmit {
		fmt.Println("✅ Submitted successfully!")
	} else {
		fmt.Println("✅ Run complete.")
	}

	return nil
}
