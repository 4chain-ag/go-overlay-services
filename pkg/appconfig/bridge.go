package appconfig

import (
	"fmt"

	"github.com/4chain-ag/go-overlay-services/pkg/internal/config"
)

// DefaultConfigFilePath is the default path to the configuration file.
const DefaultConfigFilePath = config.DefaultConfigFilePath

// NewLoader creates a new configuration loader with the given environment prefix.
func NewLoader(envPrefix string) *config.Loader[Config] {
	return config.NewLoader(Defaults, envPrefix)
}

// ToYAMLFile writes the configuration to a YAML file at the given path.
func ToYAMLFile(cfg *Config, path string) error {
	if err := config.ToYAMLFile(cfg, path); err != nil {
		return fmt.Errorf("failed to write YAML file: %w", err)
	}
	return nil
}

// ToJSONFile writes the configuration to a JSON file at the given path.
func ToJSONFile(cfg *Config, path string) error {
	if err := config.ToJSONFile(cfg, path); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}
	return nil
}

// ToEnvFile writes the configuration to an environment file at the given path.
func ToEnvFile(cfg *Config, path string) error {
	if err := config.ToEnvFile(cfg, path, cfg.AppName); err != nil {
		return fmt.Errorf("failed to write env file: %w", err)
	}
	return nil
}

// SupportedExts returns the list of supported configuration file extensions.
func SupportedExts() []string {
	return config.SupportedExts
}
