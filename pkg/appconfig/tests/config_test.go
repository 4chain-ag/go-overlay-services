package config_test

import (
	"testing"

	config "github.com/4chain-ag/go-overlay-services/pkg/appconfig"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestLoad_ShouldApplyAllDefaults_WhenNoConfigFileExists(t *testing.T) {
	// Given
	loader := config.NewLoader("OVERLAY")

	// When
	actual, err := loader.Load()
	expected := config.Defaults()
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
		Mongo: config.MongoDB{
			URI:      "mongodb://192.168.0.1:27017",
			Database: "mydb",
			Username: "admin",
			Password: "admin",
			AuthDB:   "admin",
		},
	}

	// Then
	require.NoError(t, err)
	require.Equal(t, expected, &actual)
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
