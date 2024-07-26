package utils

import (
	"os"
	"time"
)

// stores all config of the application
type Config struct {
	PostgresUser        string
	PostgresPassword    string
	PostgresDB          string
	PostgresContainer   string
	PostgresImage       string
	DBSource            string
	SchemaDir           string
	TokenSymmetricKey   string
	AccessTokenDuration time.Duration
	Environment         string
	ServerPort          string
}

func LoadConfig() (*Config, error) {
	err := LoadEnv()
	if err != nil {
		return nil, err
	}

	accessTokenDuration, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		PostgresUser:        os.Getenv("POSTGRES_USER"),
		PostgresPassword:    os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:          os.Getenv("POSTGRES_DB"),
		PostgresContainer:   os.Getenv("POSTGRES_CONTAINER_NAME"),
		PostgresImage:       os.Getenv("POSTGRES_IMAGE"),
		DBSource:            os.Getenv("DB_SOURCE"),
		SchemaDir:           os.Getenv("SCHEMA_DIR"),
		TokenSymmetricKey:   os.Getenv("TOKEN_SYMMETRIC_KEY"),
		AccessTokenDuration: accessTokenDuration,
		Environment:         os.Getenv("GIN_MODE"),
		ServerPort:          os.Getenv("PORT"),
	}

	return cfg, nil
}
