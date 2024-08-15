package lib

import (
	"github.com/joho/godotenv"
	"os"
)

type ServicesCfg struct {
	RabbitUrl string `validate:"required"`
	RedisUrl  string `validate:"required"`
	RedisPort string `validate:"required"`
}

// LoadEnv reads environment variables from a .env file and returns them as a map.
func LoadEnv() map[string]string {
	// Read environment variables from a .env file.
	config, err := godotenv.Read()
	if err != nil {
		LogError(err)
		os.Exit(1)
	}
	return config
}

func GetServicesCfg() ServicesCfg {
	env := LoadEnv()
	return ServicesCfg{
		RabbitUrl: env["RABBITMQ_URL"],
		RedisUrl:  env["REDIS_URL"],
		RedisPort: env["REDIS_PORT"],
	}
}
