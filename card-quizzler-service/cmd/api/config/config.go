package config

import (
	lib "github.com/Salladin95/card-quizzler-microservices/shared"
	"github.com/go-playground/validator/v10"
)

// Config holds the complete configuration for the application.
type Config struct {
	GrpcPort string `validate:"required"`
	DbUrl    string `validate:"required"`
	lib.ServicesCfg
}

// NewConfig creates a new configuration instance by loading environment variables and validating them.
func NewConfig() (*Config, error) {
	// Load environment variables from a .env file.
	env := lib.LoadEnv()
	lib.LogInfo("CARD QUIZZLER SERVICE ENV", env)

	// Create an AppCfg instance from the loaded environment variables.
	appCfg := Config{
		GrpcPort:    env["GRPC_PORT"],
		DbUrl:       env["DB_URL"],
		ServicesCfg: lib.GetServicesCfg(env),
	}

	// Validate the AppCfg structure using the validator package.
	validate := validator.New()

	if err := validate.Struct(appCfg); err != nil {
		return nil, err
	}

	return &appCfg, nil
}
