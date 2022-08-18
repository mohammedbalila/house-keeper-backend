package config

import (
	"fmt"

	"github.com/mustafabalila/golang-api/utils/validator"

	"github.com/caarlos0/env/v6"
	v "github.com/go-playground/validator/v10"
)

// config is the application configuration
// read from ENV vars
type config struct {
	ServiceName string `env:"SERVICE_NAME,required"`
	AppEnv      string `env:"APP_ENV"     envDefault:"development"`
	HOST        string `env:"HOST"        envDefault:"localhost"`
	PORT        int    `env:"PORT"        envDefault:"3000"`
	LogLevel    string `env:"LOG_LEVEL"   envDefault:"info"`
	DatabaseUrl string `env:"DATABASE_URL,required"`
	JWTSecret   string `env:"JWT_SECRET,required"`
	FCMAPIKey   string `env:"FCM_API_KEY,required"`
}

// parse parses, validates and then returns the application
// configuration based on ENV variables
func parse(val *v.Validate) (*config, error) {
	cfg := &config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse env: %w", err)
	}

	if err := val.Struct(cfg); err != nil {
		return nil, fmt.Errorf("failed to project env on struct: %w", err)
	}

	return cfg, nil
}

// GetConfig returns the current config
func GetConfig() *config {
	validate := validator.New()
	config, err := parse(validate)
	if err != nil {
		panic(err)
	}
	return config
}
