package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// getUserID gets the user ID from the config file
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
