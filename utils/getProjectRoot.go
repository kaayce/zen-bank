package utils

import (
	"os"
	"path/filepath"
)

const envFileName = "app.env"

func GetProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, envFileName)); !os.IsNotExist(err) {
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
