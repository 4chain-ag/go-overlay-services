package appconfig

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// MongoDB is the configuration struct for MongoDB connections.
type MongoDB struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	AuthDB   string `mapstructure:"auth_db"`
}

// DefaultMongo provides the default MongoDB configuration.
func DefaultMongo() MongoDB {
	return MongoDB{
		URI:      "mongodb://localhost:27017",
		Database: "overlay",
		Username: "",
		Password: "",
		AuthDB:   "admin",
	}
}

// validate performs validation on the MongoDB configuration.
func (cfg *MongoDB) validate() error {
	if strings.TrimSpace(cfg.URI) == "" {
		return errors.New("MongoDB URI must not be empty")
	}
	if _, err := url.ParseRequestURI(cfg.URI); err != nil {
		return fmt.Errorf("invalid MongoDB URI: %w", err)
	}
	if strings.TrimSpace(cfg.Database) == "" {
		return errors.New("MongoDB database name must not be empty")
	}
	return nil
}

// HasCredentials returns true if both username and password are set.
func (cfg *MongoDB) HasCredentials() bool {
	return strings.TrimSpace(cfg.Username) != "" && strings.TrimSpace(cfg.Password) != ""
}
