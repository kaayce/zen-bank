package utils

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	var envPath string

	env := os.Getenv("GIN_MODE")

	switch env {
	case "test":
		envPath = "example.env"
	case "debug", "release":
		envPath = ".env"
	default:
		envPath = "example.env"
	}

	rootPath, err := GetProjectRoot()
	if err != nil {
		return err
	}

	envPath = filepath.Join(rootPath, envPath)

	log.Printf("Loading env file... GIN_MODE: %v, EnvPath: %v\n", env, envPath)

	// Load the environment variables from the determined file
	err = godotenv.Load(envPath)
	if err != nil {
		return err
	}

	return nil
}
