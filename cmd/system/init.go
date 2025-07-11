package system

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/IshaanNene/EliteCode-brew/internal/utils"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Elitecode environment",
	RunE:  runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	spinner := utils.NewSpinner("Initializing Elitecode environment...")
	spinner.Start()

	// Get home directory
	home, err := homedir.Dir()
	if err != nil {
		spinner.Stop()
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	// Create .elitecode directory
	elitecodeDir := filepath.Join(home, ".elitecode")
	if err := os.MkdirAll(elitecodeDir, 0755); err != nil {
		spinner.Stop()
		return fmt.Errorf("failed to create .elitecode directory: %v", err)
	}

	// Create subdirectories
	subdirs := []string{"cache", "problems", "templates", "logs"}
	for _, subdir := range subdirs {
		subdirPath := filepath.Join(elitecodeDir, subdir)
		if err := os.MkdirAll(subdirPath, 0755); err != nil {
			spinner.Stop()
			return fmt.Errorf("failed to create %s directory: %v", subdir, err)
		}
	}

	// Create default config file
	configPath := filepath.Join(elitecodeDir, "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := `{
	"api_url": "https://api.elitecode.dev",
	"editor": "code",
	"default_language": "c",
	"auto_open_editor": true,
	"show_execution_time": true
}`
		if err := os.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
			spinner.Stop()
			return fmt.Errorf("failed to create config file: %v", err)
		}
	}

	// Create Docker templates
	if err := createDockerTemplates(elitecodeDir); err != nil {
		spinner.Stop()
		return fmt.Errorf("failed to create Docker templates: %v", err)
	}

	spinner.Stop()
	color.Green("✅ Elitecode environment initialized successfully!")
	color.Blue("📁 Configuration directory: %s", elitecodeDir)
	color.Blue("📝 Config file: %s", configPath)
	
	return nil
}

func createDockerTemplates(elitecodeDir string) error {
	templatesDir := filepath.Join(elitecodeDir, "templates")
	
	templates := map[string]string{
		"c.dockerfile": `FROM gcc:latest
WORKDIR /app
COPY . .
RUN gcc -o solution main.c
CMD ["./solution"]`,
		"cpp.dockerfile": `FROM gcc:latest
WORKDIR /app
COPY . .
RUN g++ -o solution main.cpp
CMD ["./solution"]`,
		"python.dockerfile": `FROM python:3.9-slim
WORKDIR /app
COPY . .
CMD ["python", "main.py"]`,
		"java.dockerfile": `FROM openjdk:11-jre-slim
WORKDIR /app
COPY . .
RUN javac Main.java
CMD ["java", "Main"]`,
		"javascript.dockerfile": `FROM node:16-alpine
WORKDIR /app
COPY . .
CMD ["node", "main.js"]`,
	}
	
	for filename, content := range templates {
		templatePath := filepath.Join(templatesDir, filename)
		if err := os.WriteFile(templatePath, []byte(content), 0644); err != nil {
			return err
		}
	}
	
	return nil
}