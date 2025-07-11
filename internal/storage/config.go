package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	APIBaseURL     string            `json:"apiBaseUrl"`
	AuthToken      string            `json:"authToken"`
	User           UserConfig        `json:"user"`
	GitHub         GitHubConfig      `json:"github"`
	Docker         DockerConfig      `json:"docker"`
	Editor         string            `json:"editor"`
	Language       string            `json:"language"`
	Debug          bool              `json:"debug"`
	CacheEnabled   bool              `json:"cacheEnabled"`
	CacheDirectory string            `json:"cacheDirectory"`
	ProblemsPath   string            `json:"problemsPath"`
	Settings       map[string]string `json:"settings"`
	LastUpdated    time.Time         `json:"lastUpdated"`
}

type UserConfig struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type GitHubConfig struct {
	Token        string `json:"token"`
	Username     string `json:"username"`
	Repository   string `json:"repository"`
	AutoPush     bool   `json:"autoPush"`
	CommitPrefix string `json:"commitPrefix"`
}

type DockerConfig struct {
	Enabled    bool   `json:"enabled"`
	Network    string `json:"network"`
	Timeout    int    `json:"timeout"`
	MemoryLimit string `json:"memoryLimit"`
	CPULimit   string `json:"cpuLimit"`
}

var (
	configDir  string
	configFile string
	config     *Config
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("Failed to get home directory: %v", err))
	}

	configDir = filepath.Join(homeDir, ".elitecode")
	configFile = filepath.Join(configDir, "config.json")
	
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create config directory: %v", err))
	}
}

func GetConfigDir() string {
	return configDir
}

func GetConfig() *Config {
	if config == nil {
		config = loadConfig()
	}
	return config
}

func loadConfig() *Config {
	// Default configuration
	defaultConfig := &Config{
		APIBaseURL:     "http://localhost:3000/api",
		AuthToken:      "",
		User:           UserConfig{},
		GitHub:         GitHubConfig{},
		Docker: DockerConfig{
			Enabled:     true,
			Network:     "elitecode_network",
			Timeout:     30,
			MemoryLimit: "512m",
			CPULimit:    "1.0",
		},
		Editor:         "code",
		Language:       "c",
		Debug:          false,
		CacheEnabled:   true,
		CacheDirectory: filepath.Join(configDir, "cache"),
		ProblemsPath:   filepath.Join(configDir, "problems"),
		Settings:       make(map[string]string),
		LastUpdated:    time.Now(),
	}

	// Check if config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Create default config file
		SaveConfig(defaultConfig)
		return defaultConfig
	}

	// Read existing config
	data, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Warning: Failed to read config file, using defaults: %v\n", err)
		return defaultConfig
	}

	var existingConfig Config
	if err := json.Unmarshal(data, &existingConfig); err != nil {
		fmt.Printf("Warning: Failed to parse config file, using defaults: %v\n", err)
		return defaultConfig
	}

	// Merge with defaults (in case new fields were added)
	mergeConfig(defaultConfig, &existingConfig)

	return &existingConfig
}

func mergeConfig(defaultConfig, existingConfig *Config) {
	// Update fields that might be missing in existing config
	if existingConfig.APIBaseURL == "" {
		existingConfig.APIBaseURL = defaultConfig.APIBaseURL
	}
	if existingConfig.Editor == "" {
		existingConfig.Editor = defaultConfig.Editor
	}
	if existingConfig.Language == "" {
		existingConfig.Language = defaultConfig.Language
	}
	if existingConfig.CacheDirectory == "" {
		existingConfig.CacheDirectory = defaultConfig.CacheDirectory
	}
	if existingConfig.ProblemsPath == "" {
		existingConfig.ProblemsPath = defaultConfig.ProblemsPath
	}
	if existingConfig.Docker.Network == "" {
		existingConfig.Docker = defaultConfig.Docker
	}
	if existingConfig.Settings == nil {
		existingConfig.Settings = make(map[string]string)
	}
}

func SaveConfig(cfg *Config) error {
	cfg.LastUpdated = time.Now()
	
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	config = cfg
	return nil
}

func ResetConfig() error {
	// Remove config file
	if err := os.Remove(configFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove config file: %w", err)
	}

	// Reset in-memory config
	config = nil

	return nil
}

func GetSetting(key string) string {
	cfg := GetConfig()
	if value, exists := cfg.Settings[key]; exists {
		return value
	}
	return ""
}

func SetSetting(key, value string) error {
	cfg := GetConfig()
	cfg.Settings[key] = value
	return SaveConfig(cfg)
}

func CreateDirectory(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}
	return nil
}

func EnsureDirectories() error {
	cfg := GetConfig()
	
	directories := []string{
		cfg.CacheDirectory,
		cfg.ProblemsPath,
		filepath.Join(cfg.CacheDirectory, "problems"),
		filepath.Join(cfg.CacheDirectory, "submissions"),
		filepath.Join(cfg.CacheDirectory, "templates"),
	}

	for _, dir := range directories {
		if err := CreateDirectory(dir); err != nil {
			return err
		}
	}

	return nil
}