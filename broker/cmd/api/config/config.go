package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"os"
)

// AppCfg represents the configuration settings for the application.
type AppCfg struct {
	BROKER_SERVICE_PORT string `validate:"required"` // Port for the broker service
	AUTH_SERVICE_PORT   string `validate:"required"` // Port for the auth service
	AUTH_SERVICE_URL    string `validate:"required"` // URL for the auth service
	RABBIT_URL          string `validate:"required"` // URL for RabbitMQ
}

// Config holds the complete configuration for the application.
type Config struct {
	AppCfg AppCfg // Application configuration settings
}

// NewConfig creates a new configuration instance by loading environment variables and validating them.
func NewConfig() (*Config, error) {
	// Load environment variables from a .env file.
	env := loadEnv()
	log.Infof("raw env - %v\n", env)

	// Create an AppCfg instance from the loaded environment variables.
	appCfg := AppCfg{
		BROKER_SERVICE_PORT: env["BROKER_SERVICE_PORT"],
		AUTH_SERVICE_PORT:   env["AUTH_SERVICE_PORT"],
		AUTH_SERVICE_URL:    env["AUTH_SERVICE_URL"],
		RABBIT_URL:          env["RABBITMQ_URL"],
	}

	// Validate the AppCfg structure using the validator package.
	validate := validator.New()
	if err := validate.Struct(appCfg); err != nil {
		return nil, err
	}

	// Create a new Config instance with the validated AppCfg.
	return &Config{
		AppCfg: appCfg,
	}, nil
}

// loadEnv reads environment variables from a .env file and returns them as a map.
func loadEnv() map[string]string {
	// Read environment variables from a .env file.
	config, err := godotenv.Read()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	return config
}
