package config_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/internal/config"
	"github.com/stretchr/testify/require"
)

func TestToEnvFile(t *testing.T) {
	// Given
	tmpDir := t.TempDir()
	configFilePath := fmt.Sprintf("%s/exported_config.env", tmpDir)
	cfg := Defaults()

	// When
	err := config.ToEnvFile(cfg, configFilePath, "TEST")

	// Then
	require.NoError(t, err)

	data, err := os.ReadFile(configFilePath)
	require.NoError(t, err)

	content := string(data)

	require.Contains(t, content, `TEST_A="default_hello"`)
	require.Contains(t, content, `TEST_B_WITH_LONG_NAME="1"`)
	require.Contains(t, content, `TEST_C_SUB_CONFIG_D_NESTED_FIELD="default_world"`)
}

func TestToEnvFile_WithEmptyPrefix(t *testing.T) {
	// Given
	tmpDir := t.TempDir()
	configFilePath := fmt.Sprintf("%s/exported_config.env", tmpDir)
	cfg := Defaults()

	// When
	err := config.ToEnvFile(cfg, configFilePath, "")

	// Then
	require.NoError(t, err)

	data, err := os.ReadFile(configFilePath)
	require.NoError(t, err)

	content := string(data)

	require.Contains(t, content, `A="default_hello"`)
	require.Contains(t, content, `B_WITH_LONG_NAME="1"`)
	require.NotContains(t, content, `_C_SUB_CONFIG_D_NESTED_FIELD="default_world"`)
}
