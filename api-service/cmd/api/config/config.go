package config

import (
	lib "github.com/Salladin95/card-quizzler-microservices/shared"
	"github.com/go-playground/validator/v10"
	"time"
)

// AppCfg represents the configuration settings for the application.
type AppCfg struct {
	ApiServicePort     string `validate:"required"`
	UserServiceUrl     string `validate:"required"`
	CardQuizServiceUrl string `validate:"required"`
}

type JwtCfg struct {
	AccessTokenExpTime  time.Duration `validate:"required"`
	JWTAccessSecret     string        `validate:"required"`
	RefreshTokenExpTime time.Duration `validate:"required"`
	JWTRefreshSecret    string        `validate:"required"`
}

// Config holds the complete configuration for the application.
type Config struct {
	AppCfg      AppCfg
	JwtCfg      JwtCfg
	ServicesCfg lib.ServicesCfg
}

// NewConfig creates a new configuration instance by loading environment variables and validating them.
func NewConfig() (*Config, error) {
	// Load environment variables from a .env file.
	env := lib.LoadEnv()
	lib.LogInfo("API SERVICE ENV", env)

	// Create an AppCfg instance from the loaded environment variables.
	appCfg := AppCfg{
		ApiServicePort:     env["API_SERVICE_PORT"],
		UserServiceUrl:     env["USER_SERVICE_URL"],
		CardQuizServiceUrl: env["CARD_QUIZ_SERVICE_URL"],
	}

	accessTokenExpireTime := lib.ParseDuration(env, "JWT_ACCESS_TOKEN_EXP", time.Hour*48)
	refreshTokenExpireTime := lib.ParseDuration(env, "JWT_REFRESH_TOKEN_EXP", time.Hour*72)

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
		AppCfg:      appCfg,
		JwtCfg:      jwtCfg,
		ServicesCfg: lib.GetServicesCfg(env),
	}, nil
}
