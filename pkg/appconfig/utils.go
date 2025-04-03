package appconfig

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

func PrettyPrint(cfg any) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config for printing: %w", err)
	}

	fmt.Println("Loaded Configuration:\n" + string(data))
	return nil
}

func PrettyPrintJSON(cfg any) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	fmt.Println("Loaded Configuration (JSON):\n" + string(data))
	return nil
}

func PrettyPrintAs(cfg any, format string) error {
	switch strings.ToLower(format) {
	case "json":
		return PrettyPrintJSON(cfg)
	case "yaml", "yml":
		return PrettyPrint(cfg)
	default:
		return fmt.Errorf("unsupported print format: %s", format)
	}
}
