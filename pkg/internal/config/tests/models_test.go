package config_test

type MockConfig struct {
	A string    `mapstructure:"a"`
	B int       `mapstructure:"b_with_long_name"`
	C SubConfig `mapstructure:"c_sub_config"`
}

type SubConfig struct {
	D string `mapstructure:"d_nested_field"`
}

func Defaults() MockConfig {
	return MockConfig{
		A: "default_hello",
		B: 1,
		C: SubConfig{
			D: "default_world",
		},
	}
}

const yamlConfig = `
b_with_long_name: 3
c_sub_config:
  d_nested_field: file_world
`

const dotEnvConfig = `
TEST_B_WITH_LONG_NAME=4
TEST_C_SUB_CONFIG_D_NESTED_FIELD="dotenv_world"
`

const dotEnvConfigEmptyPrefix = `
B_WITH_LONG_NAME=4
C_SUB_CONFIG_D_NESTED_FIELD="dotenv_world"
`

const jsonConfig = `
{
	"b_with_long_name": 5,
	"c_sub_config": {
		"d_nested_field": "json_world"
	}
}
`
