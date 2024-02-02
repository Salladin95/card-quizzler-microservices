package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"log"
	"os"
)

// AppCfg represents the configuration settings for the application.
type AppCfg struct {
	GrpcPort  string `validate:"required"` // Port for the auth service
	RabbitUrl string `validate:"required"` // URL for RabbitMQ
}

// FireBaseCfg represents the configuration settings for the firebase.
type FireBaseCfg struct {
	FireBaseAccKey    string `validate:"required"` // PATH TO FIRE BASE ACC KEY FILE
	FireBaseProjectId string `validate:"required"` // IF OF FIREBASE PROJECT
}

// Config holds the complete configuration for the application.
type Config struct {
	AppCfg      AppCfg      // Application configuration settings
	FireBaseCfg FireBaseCfg // Firebase configuration settings
}

// NewConfig creates a new configuration instance by loading environment variables and validating them.
func NewConfig() (*Config, error) {
	// Load environment variables from a .env file.
	env := loadEnv()

	// Create an AppCfg instance from the loaded environment variables.
	appCfg := AppCfg{
		GrpcPort:  env["GRPC_PORT"],
		RabbitUrl: env["RABBITMQ_URL"],
	}

	// Create an AppCfg instance from the loaded environment variables.
	fireBaseCfg := FireBaseCfg{
		FireBaseAccKey:    env["FIRE_BASE_ACC_KEY"],
		FireBaseProjectId: env["FIRE_BASE_PROJECT_ID"],
	}

	// Validate the AppCfg structure using the validator package.
	validate := validator.New()
	if err := validate.Struct(appCfg); err != nil {
		return nil, err
	}

	// Create a new Config instance with the validated AppCfg.
	return &Config{
		AppCfg:      appCfg,
		FireBaseCfg: fireBaseCfg,
	}, nil
}

// loadEnv reads environment variables from a .env file and returns them as a map.
func loadEnv() map[string]string {
	// Read environment variables from a .env file.
	config, err := godotenv.Read()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	return config
}
