package loaders_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server/config/loaders"
	"github.com/stretchr/testify/require"
)

func TestDefaults(t *testing.T) {
	// given:
	l := loaders.NewLoader(NewExporterTestConfig, "TEST")

	// when:
	cfg, err := l.Load()

	// then:
	require.NoError(t, err)
	require.Equal(t, "default_hello", cfg.A)
	require.Equal(t, 1, cfg.B)
	require.Equal(t, "default_world", cfg.C.D)
}

func TestEnvVariables(t *testing.T) {
	// given:
	l := loaders.NewLoader(NewExporterTestConfig, "TEST")

	// and:
	t.Setenv("TEST_B_WITH_LONG_NAME", "2")
	t.Setenv("TEST_C_SUB_CONFIG_D_NESTED_FIELD", "env_world")

	// when:
	cfg, err := l.Load()

	// then:
	require.NoError(t, err)
	require.Equal(t, "default_hello", cfg.A)
	require.Equal(t, 2, cfg.B)
	require.Equal(t, "env_world", cfg.C.D)
}

func TestFileConfig(t *testing.T) {
	// given:
	l := loaders.NewLoader(NewExporterTestConfig, "TEST")

	// and:
	configFilePath := tempConfig(t, yamlConfig, "yaml")

	// when:
	err := l.SetConfigFilePath(configFilePath)

	// then:
	require.NoError(t, err)

	// and:
	cfg, err := l.Load()

	// then:
	require.NoError(t, err)
	require.Equal(t, "default_hello", cfg.A)
	require.Equal(t, 3, cfg.B)
	require.Equal(t, "file_world", cfg.C.D)
}

func TestDotEnvConfig(t *testing.T) {
	// given:
	l := loaders.NewLoader(NewExporterTestConfig, "TEST")

	// and:
	t.Setenv("TEST_A", "env_hello")

	// and:
	configFilePath := tempConfig(t, dotEnvConfig, "env")

	// when:
	err := l.SetConfigFilePath(configFilePath)

	// then:
	require.NoError(t, err)

	// and:
	cfg, err := l.Load()

	// then:
	require.NoError(t, err)
	require.Equal(t, "env_hello", cfg.A)
	require.Equal(t, 4, cfg.B)
	require.Equal(t, "dotenv_world", cfg.C.D)
}

func TestJSONConfig(t *testing.T) {
	// given:
	l := loaders.NewLoader(NewExporterTestConfig, "TEST")

	// and:
	t.Setenv("TEST_A", "env_hello")

	// and:
	configFilePath := tempConfig(t, jsonConfig, "json")

	// when:
	err := l.SetConfigFilePath(configFilePath)

	// then:
	require.NoError(t, err)

	// and:
	cfg, err := l.Load()

	// then:
	require.NoError(t, err)
	require.Equal(t, "env_hello", cfg.A)
	require.Equal(t, 5, cfg.B)
	require.Equal(t, "json_world", cfg.C.D)
}

func TestMixedConfig(t *testing.T) {
	// given:
	l := loaders.NewLoader(NewExporterTestConfig, "TEST")

	// and:
	t.Setenv("TEST_B_WITH_LONG_NAME", "2")

	// and:
	configFilePath := tempConfig(t, yamlConfig, "yaml")

	// when:
	err := l.SetConfigFilePath(configFilePath)

	// then:
	require.NoError(t, err)

	// and:
	cfg, err := l.Load()

	// then:
	require.NoError(t, err)
	require.Equal(t, "default_hello", cfg.A)
	require.Equal(t, 2, cfg.B)
	require.Equal(t, "file_world", cfg.C.D)
}

func TestWithEmptyPrefix(t *testing.T) {
	// given:
	l := loaders.NewLoader(NewExporterTestConfig, "")

	// and:
	t.Setenv("A", "env_hello")

	// and:
	configFilePath := tempConfig(t, dotEnvConfigEmptyPrefix, "env")

	// when:
	err := l.SetConfigFilePath(configFilePath)

	// then:
	require.NoError(t, err)

	// and:
	cfg, err := l.Load()

	// then:
	require.NoError(t, err)
	require.Equal(t, "env_hello", cfg.A)
	require.Equal(t, 4, cfg.B)
	require.Equal(t, "dotenv_world", cfg.C.D)
}

func TestEnvOverridesDotEnv(t *testing.T) {
	// given:
	l := loaders.NewLoader(NewExporterTestConfig, "TEST")

	// and:
	t.Setenv("TEST_B_WITH_LONG_NAME", "2")
	t.Setenv("TEST_C_SUB_CONFIG_D_NESTED_FIELD", "env_world")

	// and:
	configFilePath := tempConfig(t, dotEnvConfig, "env")

	// when:
	err := l.SetConfigFilePath(configFilePath)

	// then:
	require.NoError(t, err)

	// and:
	cfg, err := l.Load()

	// then:
	require.NoError(t, err)
	require.Equal(t, "default_hello", cfg.A)
	require.Equal(t, 2, cfg.B)
	require.Equal(t, "env_world", cfg.C.D)
}

func tempConfig(t *testing.T, content, extension string) string {
	tmpDir := t.TempDir()
	configFilePath := fmt.Sprintf("%s/config.%s", tmpDir, extension)
	err := os.WriteFile(configFilePath, []byte(content), 0644)
	require.NoError(t, err)

	return configFilePath
}
