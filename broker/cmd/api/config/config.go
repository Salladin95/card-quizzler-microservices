package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"os"
	"time"
)

// AppCfg represents the configuration settings for the application.
type AppCfg struct {
	BrokerServicePort string `validate:"required"`
	AuthServiceUrl    string `validate:"required"`
	RabbitUrl         string `validate:"required"`
	RedisUrl          string `validate:"required"`
	RedisPort         string `validate:"required"`
}

type JwtCfg struct {
	AccessTokenExpTime  time.Duration `validate:"required"`
	JWTAccessSecret     string        `validate:"required"`
	RefreshTokenExpTime time.Duration `validate:"required"`
	JWTRefreshSecret    string        `validate:"required"`
}

// Config holds the complete configuration for the application.
type Config struct {
	AppCfg AppCfg
	JwtCfg JwtCfg
}

// NewConfig creates a new configuration instance by loading environment variables and validating them.
func NewConfig() (*Config, error) {
	// Load environment variables from a .env file.
	env := loadEnv()

	// Create an AppCfg instance from the loaded environment variables.
	appCfg := AppCfg{
		BrokerServicePort: env["BROKER_SERVICE_PORT"],
		AuthServiceUrl:    env["AUTH_SERVICE_URL"],
		RabbitUrl:         env["RABBITMQ_URL"],
		RedisUrl:          env["REDIS_URL"],
		RedisPort:         env["REDIS_PORT"],
	}

	accessTokenExpireTime := parseDuration(env, "JWT_ACCESS_TOKEN_EXP", time.Hour*48)
	refreshTokenExpireTime := parseDuration(env, "JWT_REFRESH_TOKEN_EXP", time.Hour*72)

	jwtCfg := JwtCfg{
		AccessTokenExpTime:  accessTokenExpireTime,
		RefreshTokenExpTime: refreshTokenExpireTime,
		JWTAccessSecret:     env["JWT_ACCESS_SECRET"],
		JWTRefreshSecret:    env["JWT_REFRESH_SECRET"],
	}

	// Validate the AppCfg structure using the validator package.
	validate := validator.New()

	if err := validate.Struct(appCfg); err != nil {
		return nil, err
	}
	if err := validate.Struct(jwtCfg); err != nil {
		return nil, err
	}

	// Create a new Config instance with the validated AppCfg.
	return &Config{
		AppCfg: appCfg,
		JwtCfg: jwtCfg,
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

// Parse duration parses hours
func parseDuration(config map[string]string, key string, defaultValue time.Duration) time.Duration {
	duration, err := time.ParseDuration(config[key] + "h")
	if err != nil {
		return defaultValue
	}
	return duration
}
