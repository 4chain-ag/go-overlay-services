package config

import (
	"fmt"

	"github.com/google/uuid"
)

// Config represents the application configuration.
type Config struct {
	AppName          string `mapstructure:"app_name"`
	Port             int    `mapstructure:"port"`
	Addr             string `mapstructure:"addr"`
	ServerHeader     string `mapstructure:"server_header"`
	AdminBearerToken string `mapstructure:"admin_bearer_token"`
	Mongo			struct {
		URI      string `mapstructure:"uri"`
		Database string `mapstructure:"database"`
	} `mapstructure:"mongo"`
}

// Defaults returns the default configuration values.
func Defaults() Config {
	return Config{
		AppName:          "Overlay API v0.0.0",
		Port:             3000,
		Addr:             "localhost",
		ServerHeader:     "Overlay API",
		AdminBearerToken: uuid.NewString(),
		Mongo: struct {
			URI      string `mapstructure:"uri"`
			Database string `mapstructure:"database"`
		}{
			URI:      "mongodb://localhost:27017",
			Database: "overlay",
		},
	}
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.AdminBearerToken == "" {
		return fmt.Errorf("admin bearer token is required")
	}
	return nil
}
