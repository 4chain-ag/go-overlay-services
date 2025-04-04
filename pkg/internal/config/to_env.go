package config

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/go-viper/mapstructure/v2"
)

func ToEnvFile(cfg any, filename string, envPrefix string) error {
	var decoded map[string]any
	if err := mapstructure.Decode(cfg, &decoded); err != nil {
		return fmt.Errorf("failed to decode config to map: %w", err)
	}

	if len(decoded) == 0 {
		return fmt.Errorf("config appears empty or unsupported, nothing to write")
	}

	flat := make(map[string]string)
	flattenMap(strings.ToUpper(envPrefix), decoded, flat)

	lines := make([]string, 0, len(flat))
	for k, v := range flat {
		lines = append(lines, fmt.Sprintf(`%s="%s"`, k, v))
	}
	sort.Strings(lines)

	content := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(filename, []byte(content), 0600); err != nil {
		return fmt.Errorf("failed to write env to file: %w", err)
	}

	return nil
}

func flattenMap(prefix string, input map[string]any, out map[string]string) {
	for k, v := range input {
		key := strings.ToUpper(k)
		if prefix != "" {
			key = prefix + "_" + key
		}

		switch val := v.(type) {
		case map[string]any:
			flattenMap(key, val, out)
		default:
			out[key] = fmt.Sprintf("%v", val)
		}
	}
}
