package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"os"
)

type AppCfg struct {
	BROKER_SERVICE_PORT string `validate:"required"`
	AUTH_SERVICE_PORT   string `validate:"required"`
	AUTH_SERVICE_URL    string `validate:"required"`
	RABBIT_URL          string `validate:"required"`
}

type Config struct {
	AppCfg AppCfg
}

func NewConfig() (*Config, error) {
	env := loadEnv()
	log.Infof("raw env - %v\n", env)
	appCfg := AppCfg{BROKER_SERVICE_PORT: env["BROKER_SERVICE_PORT"], AUTH_SERVICE_PORT: env["AUTH_SERVICE_PORT"], AUTH_SERVICE_URL: env["AUTH_SERVICE_URL"], RABBIT_URL: env["RABBITMQ_URL"]}
	validate := validator.New()
	if err := validate.Struct(appCfg); err != nil {
		return nil, err
	}
	return &Config{
		AppCfg: appCfg,
	}, nil
}

func loadEnv() map[string]string {
	config, err := godotenv.Read()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	return config
}
