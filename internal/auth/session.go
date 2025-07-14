package auth

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Session struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	Token    string `json:"token,omitempty"`
}

var sessionFile = filepath.Join(os.Getenv("HOME"), ".elitecode", "session.json")

func SaveSession(s Session) error {
	os.MkdirAll(filepath.Dir(sessionFile), os.ModePerm)
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(sessionFile, data, 0644)
}

func LoadSession() (*Session, error) {
	data, err := os.ReadFile(sessionFile)
	if err != nil {
		return nil, err
	}
	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func ClearSession() error {
	return os.Remove(sessionFile)
}

func IsLoggedIn() bool {
	_, err := LoadSession()
	return err == nil
}
