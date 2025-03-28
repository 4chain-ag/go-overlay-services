package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestLoad_ShouldApplyAllDefaults_WhenNoConfigFileExists(t *testing.T) {
	// Given
	loader := config.NewLoader("OVERLAY")

	// When
	actual, err := loader.Load()
	expected := config.DefaultConfig()
	expected.AdminBearerToken = actual.AdminBearerToken

	// Then
	require.NoError(t, err)
	require.Equal(t, expected, actual)

	_, err = uuid.Parse(actual.AdminBearerToken)
	require.NoError(t, err, "admin token should be a valid UUID")
}

func TestLoad_ShouldOverrideDefaults_WhenConfigFileProvidesValues(t *testing.T) {
	// Given
	loader := config.NewLoader("OVERLAY")
	require.NoError(t, loader.SetConfigFilePath("testdata/config.yaml"))

	// When
	actual, err := loader.Load()

	expected := &config.Config{
		AppName:          "CustomApp",
		Port:             9999,
		Addr:             "127.0.0.1",
		ServerHeader:     "CustomHeader",
		AdminBearerToken: "secret-token",
	}

	// Then
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestSetConfigFilePath_ShouldReturnError_WhenUnsupportedExtension(t *testing.T) {
	// Given
	loader := config.NewLoader("OVERLAY")

	// When
	err := loader.SetConfigFilePath("config.txt")

	// Then
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported config file extension")
}

func TestExportToYAML_ShouldWriteFile_WhenConfigIsValid(t *testing.T) {
	// Given
	loader := config.NewLoader("OVERLAY")
	_, err := loader.Load()
	require.NoError(t, err)

	// When
	tmpFile := filepath.Join(t.TempDir(), "config.yaml")
	err = loader.ToYAML(tmpFile)

	// Then
	require.NoError(t, err)
	data, err := os.ReadFile(tmpFile)
	require.NoError(t, err)
	require.Contains(t, string(data), "addr: localhost")
}

func TestExportToJSON_ShouldWriteFile_WhenConfigIsValid(t *testing.T) {
	// Given
	loader := config.NewLoader("OVERLAY")
	_, err := loader.Load()
	require.NoError(t, err)

	// When
	tmpFile := filepath.Join(t.TempDir(), "config.json")
	err = loader.ToJSON(tmpFile)

	// Then
	require.NoError(t, err)
	data, err := os.ReadFile(tmpFile)
	require.NoError(t, err)
	require.Contains(t, string(data), `"Addr": "localhost"`)
}

func TestExportToEnv_ShouldWriteFlatEnvFile_WhenConfigIsValid(t *testing.T) {
	// Given
	loader := config.NewLoader("OVERLAY")
	_, err := loader.Load()
	require.NoError(t, err)

	// When
	tmpFile := filepath.Join(t.TempDir(), ".env")
	err = loader.ToEnv(tmpFile)

	// Then
	require.NoError(t, err)
	data, err := os.ReadFile(tmpFile)
	require.NoError(t, err)

	content := string(data)
	require.Contains(t, content, "ADDR=localhost")
	require.Contains(t, content, "PORT=3000")
	require.Contains(t, content, "APP_NAME=Overlay API v0.0.0")
}

func TestExportToEnv_ShouldFail_WhenFilePathIsInvalid(t *testing.T) {
	// Given
	loader := config.NewLoader("OVERLAY")
	_, _ = loader.Load()

	// When
	err := loader.ToEnv("/invalid/!!/envfile")

	// Then
	require.Error(t, err)
}
