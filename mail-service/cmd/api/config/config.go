package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"os"
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
	config := loadEnv()

	cfg := Config{
		Port:              config["PORT"],
		RabbitUrl:         config["RABBITMQ_URL"],
		SmtpServerAddress: config["SMTP_SERVER_ADDRESS"],
		SmtpAuthAddress:   config["SMTP_AUTH_ADDRESS"],
		EmailName:         config["APP_EMAIL_NAME"],
		EmailAddress:      config["APP_EMAIL"],
		EmailPassword:     config["APP_EMAIL_PASSWORD"],
	}

	// Validate the cfg structure using the validator package.
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func loadEnv() map[string]string {
	config, err := godotenv.Read()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		os.Exit(1)
	}
	return config
}
