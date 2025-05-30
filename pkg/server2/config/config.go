package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/config/exporters"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/config/loaders"
)

// Config contains configuration settings for the overlay-engine API and its dependencies.
type Config struct {
	Server server2.Config `mapstructure:"server"`
}

// Export writes the configuration to the file at the specified path.
// It formats the file content based on the file extension:
// - JSON for ".json" files
// - Environment variables for ".env" or ".dotenv" files
// - YAML for ".yaml" or ".yml" files
func (c *Config) Export(path string) error {
	ext := strings.TrimPrefix(filepath.Ext(path), ".")
	var err error
	switch ext {
	case "json":
		err = exporters.ToJSONFile(c, path)
	case "env", "dotenv":
		err = exporters.ToEnvFile(c, path, strings.Replace(c.Server.AppName, " ", "_", -1))
	default: // yaml, yml
		err = exporters.ToYAMLFile(c, path)
	}

	if err != nil {
		return fmt.Errorf("failed to export configuration: %w", err)
	}
	return nil
}

// NewDefault returns a Config with default HTTP server and MongoDB settings.
func NewDefault() Config {
	return Config{
		Server: server2.DefaultConfig,
	}
}

// LoadFromPath loads the server configuration from the specified file path.
// It initializes a new loader using the default config provider and the environment prefix.
// The function attempts to read and decode the config file, pretty-prints the configuration as JSON,
// and returns the extracted server configuration on success. An error is returned if any step fails.
func LoadFromPath(path, env string) (server2.Config, error) {
	loader := loaders.NewLoader(NewDefault, env)
	err := loader.SetConfigFilePath(path)
	if err != nil {
		return server2.Config{}, fmt.Errorf("invalid config file path: %w", err)
	}

	cfg, err := loader.Load()
	if err != nil {
		return server2.Config{}, fmt.Errorf("config loader load operation failed: %w", err)
	}

	err = PrettyPrintAs(cfg, "json")
	if err != nil {
		return server2.Config{}, fmt.Errorf("config pretty print operation failed: %w", err)
	}
	return cfg.Server, nil
}
