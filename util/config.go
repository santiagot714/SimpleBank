package util

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variables.
type Config struct {
	AppVersion          string        `mapstructure:"APP_VERSION"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBHost              string        `mapstructure:"DB_HOST"`
	DBPort              string        `mapstructure:"DB_PORT"`
	DBUser              string        `mapstructure:"DB_USER"`
	DBPassword          string        `mapstructure:"DB_PASSWORD"`
	DBName              string        `mapstructure:"DB_NAME"`
	DatabaseURL         string        `mapstructure:"DATABASE_URL"`
	TestDBHost          string        `mapstructure:"TEST_DB_HOST"`
	TestDBPort          string        `mapstructure:"TEST_DB_PORT"`
	TestDBUser          string        `mapstructure:"TEST_DB_USER"`
	TestDBPassword      string        `mapstructure:"TEST_DB_PASSWORD"`
	TestDBName          string        `mapstructure:"TEST_DB_NAME"`
	TestDatabaseURL     string        `mapstructure:"TEST_DATABASE_URL"`
	TokenKey            string        `mapstructure:"PASETO_TOKEN_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

// configKeys lists all environment variable names mapped to Config fields.
// Explicit binding is required so that viper.Unmarshal populates the struct
// from environment variables even when no config file is present (e.g. CI).
var configKeys = []string{
	"APP_VERSION", "SERVER_ADDRESS",
	"DB_DRIVER", "DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME",
	"DATABASE_URL",
	"TEST_DB_HOST", "TEST_DB_PORT", "TEST_DB_USER", "TEST_DB_PASSWORD",
	"TEST_DB_NAME", "TEST_DATABASE_URL",
	"PASETO_TOKEN_KEY", "ACCESS_TOKEN_DURATION",
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigFile(filepath.Join(path, ".env"))

	viper.AutomaticEnv()

	for _, key := range configKeys {
		if err = viper.BindEnv(key); err != nil {
			return
		}
	}

	if err = viper.ReadInConfig(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
