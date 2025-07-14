package utils

import (
	"os"
	"path/filepath"
)

func EnsureDir(path string) error {
	return os.MkdirAll(filepath.Dir(path), os.ModePerm)
}

func WriteFile(path string, content []byte) error {
	if err := EnsureDir(path); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0644)
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
