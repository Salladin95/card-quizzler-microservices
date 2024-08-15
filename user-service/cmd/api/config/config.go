package config

import (
	lib "github.com/Salladin95/card-quizzler-microservices/shared"
	"github.com/go-playground/validator/v10"
)

type AppCfg struct {
	GrpcPort string `validate:"required"`
	lib.ServicesCfg
}

type MongoCfg struct {
	MongoUrl          string `validate:"required"`
	MongoUsername     string `validate:"required"`
	MongoUserPassword string `validate:"required"`
	MongoDbName       string `validate:"required"`
}

// Config holds the complete configuration for the application.
type Config struct {
	AppCfg   AppCfg
	MongoCfg MongoCfg
}

// NewConfig creates a new configuration instance by loading environment variables and validating them.
func NewConfig() (*Config, error) {
	// Load environment variables from a .env file.
	env := lib.LoadEnv()
	lib.LogInfo("USER SERVICE ENV", env)

	// Create an AppCfg instance from the loaded environment variables.
	appCfg := AppCfg{
		GrpcPort:    env["GRPC_PORT"],
		ServicesCfg: lib.GetServicesCfg(),
	}

	// Validate the AppCfg structure using the validator package.
	validate := validator.New()
	if err := validate.Struct(appCfg); err != nil {
		return nil, err
	}

	mongoCfg := MongoCfg{
		MongoUrl:          env["MONGO_URL"],
		MongoUsername:     env["MONGO_USERNAME"],
		MongoUserPassword: env["MONGO_PASSWORD"],
		MongoDbName:       env["MONGO_DB"],
	}

	if err := validate.Struct(mongoCfg); err != nil {
		return nil, err
	}

	// Create a new Config instance with the validated AppCfg.
	return &Config{
		AppCfg:   appCfg,
		MongoCfg: mongoCfg,
	}, nil
}
