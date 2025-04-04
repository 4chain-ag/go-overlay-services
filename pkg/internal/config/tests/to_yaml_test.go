package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/internal/config"
	"github.com/stretchr/testify/require"
)

func TestToYAMLFile(t *testing.T) {
	// given:
	tmpDir := t.TempDir()
	configFilePath := fmt.Sprintf("%s/exported_config.yaml", tmpDir)

	// and:
	cfg := Defaults()

	// when:
	err := config.ToYAMLFile(cfg, configFilePath)

	// then:
	require.NoError(t, err)

	yamlFile, err := os.ReadFile(configFilePath)
	require.NoError(t, err)

	require.Contains(t, string(yamlFile), "a: default_hello")
	require.Contains(t, string(yamlFile), "b_with_long_name: 1")
	require.Contains(t, string(yamlFile), "d: default_world")
}

func TestExportToYAML_ShouldWriteFile_WhenConfigIsValid(t *testing.T) {
	// Given
	l := config.NewLoader(Defaults, "OVERLAY")
	_, err := l.Load()
	require.NoError(t, err)

	// When
	tmpFile := filepath.Join(t.TempDir(), "config.yaml")
	err = config.ToYAMLFile(Defaults(), tmpFile)

	// Then
	require.NoError(t, err)
	data, err := os.ReadFile(tmpFile)
	require.NoError(t, err)
	require.Contains(t, string(data), "a: default_hello")
	require.Contains(t, string(data), "d_nested_field: default_world")
}
