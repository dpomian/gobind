package utils

import (
	"os"
	"time"
)

type Config struct {
	DBSource            string        `mapstructure:"BINDER_DB_SOURCE"`
	DBDriver            string        `mapstructure:"BINDER_DB_DRIVER"`
	ServerAddress       string        `mapstructure:"BINDER_SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"BINDER_TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"BINDER_ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (Config, error) {
	config := Config{
		DBSource:            os.Getenv("BINDER_DB_SOURCE"),
		DBDriver:            os.Getenv("BINDER_DB_DRIVER"),
		ServerAddress:       ":5050",
		TokenSymmetricKey:   RandomString(32),
		AccessTokenDuration: time.Duration(15 * time.Minute),
	}

	return config, nil
}
