package config

import (
	lib "github.com/Salladin95/card-quizzler-microservices/shared"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	Port              string `validate:"required"`
	RabbitUrl         string `validate:"required"`
	SmtpAuthAddress   string `validate:"required"`
	SmtpServerAddress string `validate:"required"`
	EmailName         string `validate:"required"`
	EmailAddress      string `validate:"required"`
	EmailPassword     string `validate:"required"`
}

func GetConfig() (*Config, error) {
	env := lib.LoadEnv()
	lib.LogInfo("MAIL SERVICE ENV", env)

	cfg := Config{
		Port:              env["PORT"],
		RabbitUrl:         env["RABBITMQ_URL"],
		SmtpServerAddress: env["SMTP_SERVER_ADDRESS"],
		SmtpAuthAddress:   env["SMTP_AUTH_ADDRESS"],
		EmailName:         env["APP_EMAIL_NAME"],
		EmailAddress:      env["APP_EMAIL"],
		EmailPassword:     env["APP_EMAIL_PASSWORD"],
	}

	// Validate the cfg structure using the validator package.
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
