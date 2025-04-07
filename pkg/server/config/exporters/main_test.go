package exporters_test

type ExporterTestConfig struct {
	A string                `mapstructure:"a"`
	B int                   `mapstructure:"b_with_long_name"`
	C ExporterTestSubConfig `mapstructure:"c_sub_config"`
}

type ExporterTestSubConfig struct {
	D string `mapstructure:"d_nested_field"`
}

func NewExporterTestConfig() ExporterTestConfig {
	return ExporterTestConfig{
		A: "default_hello",
		B: 1,
		C: ExporterTestSubConfig{
			D: "default_world",
		},
	}
}
