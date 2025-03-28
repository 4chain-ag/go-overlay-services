package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/gookit/slog"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// DefaultConfigFilePath is the default path to the configuration file
const DefaultConfigFilePath = "config.yaml"

// Load implements the Loader and Exporter interfaces
type Load struct {
	cfg            *Config
	envPrefix      string
	configFilePath string
	configFileExt  string
	viper          *viper.Viper
}

// NewLoader creates a new configuration loader
func NewLoader(envPrefix string) *Load {
	return &Load{
		cfg:            NewConfig(),
		envPrefix:      envPrefix,
		configFilePath: DefaultConfigFilePath,
		viper:          viper.New(),
	}
}

// SetConfigFilePath sets the path to the configuration file
func (l *Load) SetConfigFilePath(path string) error {
	ext := strings.TrimPrefix(filepath.Ext(path), ".")
	if ext != "yaml" && ext != "yml" && ext != "json" {
		return fmt.Errorf("unsupported config file extension: %s", ext)
	}
	l.configFilePath = path
	l.configFileExt = ext
	return nil
}

// Load reads the configuration from the file and environment variables
func (l *Load) Load() (*Config, error) {
	l.setViperDefaults()
	l.prepareViper()

	if err := l.loadFromFile(); err != nil {
		return l.cfg, err
	}
	if err := l.viperToCfg(); err != nil {
		return l.cfg, err
	}

	slog.Info("Config loaded successfully")
	return l.cfg, nil
}

func (l *Load) setViperDefaults() {
	defaultsMap := map[string]interface{}{}
	if err := mapstructure.Decode(NewConfig(), &defaultsMap); err != nil {
		slog.Errorf("error while setting defaults: %v", err)
		return
	}
	for k, v := range defaultsMap {
		l.viper.SetDefault(k, v)
	}
}

func (l *Load) prepareViper() {
	l.viper.SetEnvPrefix(l.envPrefix)
	l.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	l.viper.AutomaticEnv()
}

func (l *Load) loadFromFile() error {
	if _, err := os.Stat(l.configFilePath); os.IsNotExist(err) {
		return nil
	}

	l.viper.SetConfigFile(l.configFilePath)
	if err := l.viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}
	slog.Infof("Loaded config from file: %s", l.configFilePath)
	return nil
}

func (l *Load) viperToCfg() error {
	if err := l.viper.Unmarshal(l.cfg); err != nil {
		return fmt.Errorf("error unmarshalling config: %w", err)
	}
	l.updateViperFromConfig()
	return nil
}

func (l *Load) updateViperFromConfig() {
	configMap := map[string]any{}
	if err := mapstructure.Decode(l.cfg, &configMap); err != nil {
		slog.Errorf("failed to encode config back to map: %v", err)
		return
	}
	for key, value := range configMap {
		l.viper.Set(key, value)
	}
}

// PrettyPrint logs the loaded config as indented JSON using slog
func (l *Load) PrettyPrint() error {
	data, err := json.MarshalIndent(l.cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	slog.Infof("Loaded Configuration:\n%s", string(data))
	return nil
}

// ToJSON writes the configuration to a JSON file
func (l *Load) ToJSON(filePath string) error {
	data, err := json.MarshalIndent(l.cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}
	slog.Infof("JSON config exported to: %s", filePath)
	return nil
}

// ToYAML writes the configuration to a YAML file
func (l *Load) ToYAML(filePath string) error {
	data, err := yaml.Marshal(l.cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}
	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write YAML file: %w", err)
	}
	slog.Infof("YAML config exported to: %s", filePath)
	return nil
}

// ToEnv writes the configuration to an environment file
func (l *Load) ToEnv(filePath string) error {
	data := map[string]any{}
	if err := mapstructure.Decode(l.cfg, &data); err != nil {
		return fmt.Errorf("failed to decode config: %w", err)
	}

	flattened := map[string]string{}
	flattenConfig("", flattened, data)

	var content string
	for key, value := range flattened {
		content += fmt.Sprintf("%s=%s\n", key, value)
	}

	if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
		return fmt.Errorf("failed to write env file: %w", err)
	}

	slog.Infof("ENV config exported to: %s", filePath)
	return nil
}
