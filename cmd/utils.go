package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/IshaanNene/EliteCode-brew/internal/models"
	"github.com/fatih/color"
)

func formatDifficulty(difficulty models.Difficulty) string {
	switch strings.ToLower(string(difficulty)) {
	case "easy":
		return color.GreenString("Easy")
	case "medium":
		return color.YellowString("Medium")
	case "hard":
		return color.RedString("Hard")
	default:
		return string(difficulty)
	}
}

func getUserID() (string, error) {
	configFile := filepath.Join(os.Getenv("HOME"), ".elitecode", "config.json")
	configBytes, err := os.ReadFile(configFile)
	if err != nil {
		return "", fmt.Errorf("error reading config file: %v", err)
	}

	var config struct {
		UID string `json:"uid"`
	}
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return "", fmt.Errorf("error parsing config: %v", err)
	}

	if config.UID == "" {
		return "", fmt.Errorf("not logged in. Please run 'elitecode login' first")
	}

	return config.UID, nil
}
