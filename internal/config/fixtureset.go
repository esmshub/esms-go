package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Fixtureset struct {
	Name             string
	OverrideSettings map[string]any `mapstructure:"override_settings" json:"override_settings" yaml:"override_settings"`
	Fixtures         []*Fixture
}

func LoadFixtureset[T any](path string) (*T, error) {
	v := viper.New()
	v.SetConfigFile(path)
	err := v.ReadInConfig() // Find and read the config file
	if err != nil {         // Handle errors reading the config file
		panic(fmt.Errorf("error loading fixtureset: %w", err))
	}

	var conf T
	v.Unmarshal(&conf)
	return &conf, nil
}
