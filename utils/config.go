package utils

import (
	"os"
	"time"
)

type Config struct {
	DBSource             string        `mapstructure:"BINDER_DB_SOURCE"`
	DBDriver             string        `mapstructure:"BINDER_DB_DRIVER"`
	ServerAddress        string        `mapstructure:"BINDER_API_SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"BINDER_TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"BINDER_ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"BINDER_REFRESH_TOKEN_DURATION"`
}

func LoadConfig(path string) (Config, error) {
	config := Config{
		DBSource:             os.Getenv("BINDER_DB_SOURCE"),
		DBDriver:             os.Getenv("BINDER_DB_DRIVER"),
		ServerAddress:        os.Getenv("BINDER_API_SERVER_ADDRESS"),
		TokenSymmetricKey:    RandomString(32),
		AccessTokenDuration:  time.Duration(15 * time.Minute),
		RefreshTokenDuration: time.Duration(24 * time.Hour),
	}

	return config, nil
}

type UIConfig struct {
	RedisUri         string `mapstructure:"REDIS_URI"`
	RedisSecret      string `mapstructure:"REDIS_SECRET"`
	ServerAddress    string `mapstructure:"BINDER_UI_SERVER_ADDRESS"`
	BinderApiBaseUrl string `mapstructure:"BINDER_API_BASE_URL"`
}

func LoadUiConfig(path string) (UIConfig, error) {
	config := UIConfig{
		RedisUri:         os.Getenv("REDIS_URI"),
		RedisSecret:      "secret",
		ServerAddress:    os.Getenv("BINDER_UI_SERVER_ADDRESS"),
		BinderApiBaseUrl: os.Getenv("BINDER_API_BASE_URL"),
	}

	return config, nil
}
