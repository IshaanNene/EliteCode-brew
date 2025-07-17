package firebase

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

// Config holds Firebase configuration
type Config struct {
	APIKey            string `json:"apiKey"`
	AuthDomain        string `json:"authDomain"`
	ProjectID         string `json:"projectId"`
	StorageBucket     string `json:"storageBucket"`
	MessagingSenderID string `json:"messagingSenderId"`
	AppID            string `json:"appId"`
	MeasurementID    string `json:"measurementId"`
}

var defaultConfig = Config{
	APIKey:            "AIzaSyDhoR2EA3uuWqWSkiiE5G1vRkeqoGHBdro",
	AuthDomain:        "elitecode-47058.firebaseapp.com",
	ProjectID:         "elitecode-47058",
	StorageBucket:     "elitecode-47058.firebasestorage.app",
	MessagingSenderID: "605884776680",
	AppID:            "1:605884776680:web:f89db0d49b021139dec95d",
	MeasurementID:    "G-4XKNK2CH3E",
}

// InitializeApp initializes the Firebase app with the provided or default configuration
func InitializeApp(ctx context.Context) (*firebase.App, error) {
	config := defaultConfig

	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting user home directory: %v", err)
	}

	// Create .elitecode directory if it doesn't exist
	configDir := filepath.Join(homeDir, ".elitecode")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("error creating config directory: %v", err)
	}

	// Create config file if it doesn't exist
	configFile := filepath.Join(configDir, "firebase_config.json")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		configJSON, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("error marshaling config: %v", err)
		}
		if err := os.WriteFile(configFile, configJSON, 0644); err != nil {
			return nil, fmt.Errorf("error writing config file: %v", err)
		}
	}

	// Initialize Firebase app
	conf := &firebase.Config{
		ProjectID:     config.ProjectID,
		StorageBucket: config.StorageBucket,
	}

	// Check for service account key file
	saKeyFile := filepath.Join(configDir, "firebase_credentials.json")
	var app *firebase.App
	if _, err := os.Stat(saKeyFile); err == nil {
		// Use service account if available
		app, err = firebase.NewApp(ctx, conf, option.WithCredentialsFile(saKeyFile))
	} else {
		// Use default credentials
		app, err = firebase.NewApp(ctx, conf)
	}

	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %v", err)
	}

	return app, nil
}

// GetConfig returns the current Firebase configuration
func GetConfig() Config {
	return defaultConfig
} 