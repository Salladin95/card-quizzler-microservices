package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"os"
	"time"
)

// Config holds the complete configuration for the application.
type Config struct {
	GrpcPort  string `validate:"required"`
	DbUrl     string `validate:"required"`
	RabbitUrl string `validate:"required"`
	RedisUrl  string `validate:"required"`
	RedisPort string `validate:"required"`
}

// NewConfig creates a new configuration instance by loading environment variables and validating them.
func NewConfig() (*Config, error) {
	// Load environment variables from a .env file.
	env := loadEnv()

	// Create an AppCfg instance from the loaded environment variables.
	appCfg := Config{
		GrpcPort:  env["GRPC_PORT"],
		DbUrl:     env["DB_URL"],
		RabbitUrl: env["RABBITMQ_URL"],
		RedisUrl:  env["REDIS_URL"],
		RedisPort: env["REDIS_PORT"],
	}

	// Validate the AppCfg structure using the validator package.
	validate := validator.New()

	if err := validate.Struct(appCfg); err != nil {
		return nil, err
	}

	return &appCfg, nil
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

// Parse duration parses hours
func parseDuration(config map[string]string, key string, defaultValue time.Duration) time.Duration {
	duration, err := time.ParseDuration(config[key] + "h")
	if err != nil {
		return defaultValue
	}
	return duration
}
