package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"os"
)

// Config holds the complete configuration for the application.
type Config struct {
	RabbitUrl string `validate:"required"` // amqp rabbit url
}

// NewConfig creates a new configuration instance by loading environment variables and validating them.
func NewConfig() (*Config, error) {
	// Load environment variables from a .env file.
	env := loadEnv()

	cfg := Config{
		RabbitUrl: env["RABBITMQ_URL"],
	}

	// Validate the AppCfg structure using the validator package.
	validate := validator.New()

	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	// Create a new Config instance with the validated AppCfg.
	return &cfg, nil
}

// loadEnv reads environment variables from a .env file and returns them as a map.
func loadEnv() map[string]string {
	// Read environment variables from a .env file.
	config, err := godotenv.Read()
	if err != nil {
		os.Exit(1)
	}

	return config
}
