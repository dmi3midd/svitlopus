package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Http represents the HTTP server configuration.
type Http struct {
	Address      string        `mapstructure:"address"`
	IdleTimeout  time.Duration `mapstructure:"idleTimeout"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
}

// Database represents the database configuration.
type Database struct {
	DbPath string `mapstructure:"dbPath"`
}

// Log
type Log struct {
	LogPath string `mapstructure:"logPath"`
}

type Config struct {
	Database Database `mapstructure:"database"`
	Http     Http     `mapstructure:"http"`
	Log      Log      `mapstructure:"log"`
}

// LoadConfig loads the configuration from the config file config.yaml.
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config into struct: %w", err)
	}

	return &cfg, nil
}
