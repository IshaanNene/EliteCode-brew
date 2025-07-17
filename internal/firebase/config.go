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

type Config struct {
	APIKey            string `json:"apiKey"`
	AuthDomain        string `json:"authDomain"`
	ProjectID         string `json:"projectId"`
	StorageBucket     string `json:"storageBucket"`
	MessagingSenderID string `json:"messagingSenderId"`
	AppID             string `json:"appId"`
	MeasurementID     string `json:"measurementId"`
}

var defaultConfig = Config{
	APIKey:            "AIzaSyCblKxoY7dXZspe-szzLN7lwYb3eFWGvXI",
	AuthDomain:        "brew-elitecode-store.firebaseapp.com",
	ProjectID:         "brew-elitecode-store",
	StorageBucket:     "brew-elitecode-store.firebasestorage.app",
	MessagingSenderID: "833944171879",
	AppID:             "1:833944171879:web:ac583b56b0c99c9b892bc2",
	MeasurementID:     "G-4XKNK2CH3E",
}

func InitializeApp(ctx context.Context) (*firebase.App, error) {
	config := defaultConfig

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting user home directory: %v", err)
	}

	configDir := filepath.Join(homeDir, ".elitecode")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("error creating config directory: %v", err)
	}

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

	conf := &firebase.Config{
		ProjectID:     config.ProjectID,
		StorageBucket: config.StorageBucket,
	}

	saKeyFile := filepath.Join(configDir, "firebase_credentials.json")
	var app *firebase.App
	if _, err := os.Stat(saKeyFile); err == nil {
		app, err = firebase.NewApp(ctx, conf, option.WithCredentialsFile(saKeyFile))
	} else {
		app, err = firebase.NewApp(ctx, conf)
	}

	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %v", err)
	}

	return app, nil
}

func GetConfig() Config {
	return defaultConfig
}
