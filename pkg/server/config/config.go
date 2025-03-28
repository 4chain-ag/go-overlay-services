package config

import "github.com/google/uuid"

// Config defines the configuration for the application.
type Config struct {
	AppName          string
	Port             int
	Addr             string
	ServerHeader     string
	AdminBearerToken string
}

// Option defines a functional option for configuring the configuration.
type Option func(*Config)

// WithAppName sets the app name for the configuration.
func WithAppName(appName string) Option {
	return func(c *Config) {
		c.AppName = appName
	}
}

// WithPort sets the port for the configuration.
func WithPort(port int) Option {
	return func(c *Config) {
		c.Port = port
	}
}

// WithAddr sets the address for the configuration.
func WithAddr(addr string) Option {
	return func(c *Config) {
		c.Addr = addr
	}
}

// WithServerHeader sets the server header for the configuration.
func WithServerHeader(serverHeader string) Option {
	return func(c *Config) {
		c.ServerHeader = serverHeader
	}
}

// WithAdminBearerToken sets the admin bearer token for the configuration.
func WithAdminBearerToken(adminBearerToken string) Option {
	return func(c *Config) {
		c.AdminBearerToken = adminBearerToken
	}
}

// NewConfig creates a new configuration instance with the specified options.
func NewConfig(opts ...Option) *Config {
	cfg := &Config{
		AppName:          "Overlay API v0.0.0",
		Port:             3000,
		Addr:             "localhost",
		ServerHeader:     "Overlay API",
		AdminBearerToken: uuid.NewString(),
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// DefaultConfig is an alias for NewConfig without options.
func DefaultConfig() *Config {
	return NewConfig()
}
