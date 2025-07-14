package utils

import (
	"errors"
	"os"
	"path/filepath"
)

func GetFirebaseCredentialsPath() (string, error) {
	path := filepath.Join("config", "serviceAccountKey.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", errors.New("Firebase credentials not found")
	}
	return path, nil
}
