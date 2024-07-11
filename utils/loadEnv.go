package utils

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	var envPath string

	env := os.Getenv("ENV")

	switch env {
	case "public", "test":
		envPath = "example.env"
	case "dev", "prod":
	default:
		envPath = ".env"
	}

	rootPath, err := GetProjectRoot()
	if err != nil {
		return err
	}

	envPath = filepath.Join(rootPath, envPath)

	// Load the environment variables from the determined file
	err = godotenv.Load(envPath)
	if err != nil {
		return err
	}

	return nil
}
