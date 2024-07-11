package utils

import (
	"os"
	"path/filepath"
)

func GetProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); !os.IsNotExist(err) {
			return dir, nil
		}

		// Check if we have reached the root directory
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			break
		}
		dir = parentDir
	}

	return "", os.ErrNotExist
}
