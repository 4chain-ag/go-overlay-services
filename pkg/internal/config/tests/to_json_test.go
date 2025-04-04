package config_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/internal/config"
	"github.com/stretchr/testify/require"
)

func TestToJSONFile(t *testing.T) {
	// given:
	tmpDir := t.TempDir()
	configFilePath := fmt.Sprintf("%s/exported_config.json", tmpDir)

	cfg := Defaults()

	// when:
	err := config.ToJSONFile(cfg, configFilePath)

	// then:
	require.NoError(t, err)

	data, err := os.ReadFile(configFilePath)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	expected := map[string]any{
		"a":                "default_hello",
		"b_with_long_name": float64(1),
		"c_sub_config": map[string]any{
			"d_nested_field": "default_world",
		},
	}

	require.Equal(t, expected, result)
}
