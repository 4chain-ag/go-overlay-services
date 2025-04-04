package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-viper/mapstructure/v2"
)

// ToJSONFile exports the given config struct into a JSON file.
func ToJSONFile(cfg any, filename string) error {
	var mapData map[string]any

	err := mapstructure.Decode(cfg, &mapData)
	if err != nil {
		return fmt.Errorf("failed to decode config to map: %w", err)
	}

	jsonData, err := json.MarshalIndent(mapData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal map to json: %w", err)
	}

	err = os.WriteFile(filename, jsonData, 0600)
	if err != nil {
		return fmt.Errorf("failed to write json to file: %w", err)
	}

	return nil
}
