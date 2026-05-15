package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Http represents the HTTP server configuration.
type Http struct {
	Address      string        `yaml:"address"`
	IdleTimeout  time.Duration `yaml:"idleTimeout"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
}

// Database represents the database configuration.
type Database struct {
	DbPath string `yaml:"dbPath"`
}

// Log
type Log struct {
	LogPath string `yaml:"logPath"`
}

// Docker
type Docker struct {
	ApiHash       string `yaml:"apiHash"`
	ApiId         string `yaml:"apiId"`
	Image         string `yaml:"image"`
	ContainerName string `yaml:"containerName"`
	Port          int    `yaml:"port"`
	Valume        string `yaml:"volume"`
	RestartAlways bool   `yaml:"restartAlways"`
}

// Bot
type Bot struct {
	Token string `yaml:"token"`
}

type Config struct {
	Database Database `yaml:"database"`
	Http     Http     `yaml:"http"`
	Log      Log      `yaml:"log"`
	Docker   Docker   `yaml:"docker"`
	Bot      Bot      `yaml:"bot"`
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
