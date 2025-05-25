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

func LoadFixtureset(path string) (*Fixtureset, error) {
	v := viper.New()
	v.SetConfigFile(path)
	err := v.ReadInConfig() // Find and read the config file
	if err != nil {         // Handle errors reading the config file
		panic(fmt.Errorf("error loading fixtureset: %w", err))
	}

	var conf Fixtureset
	v.Unmarshal(&conf)
	return &conf, validateFixtureset(&conf)
}

func validateFixtureset(fs *Fixtureset) error {
	if fs.Name == "" {
		return fmt.Errorf("fixtureset name is required")
	}

	if len(fs.Fixtures) == 0 {
		return fmt.Errorf("fixtureset must have at least one fixture")
	}

	for _, f := range fs.Fixtures {
		if f.HomeTeamCode == "" && f.HomeTeamsheet == "" {
			return fmt.Errorf("fixture must have either home_teamsheet or home_team set")
		}
		if f.AwayTeamCode == "" && f.AwayTeamsheet == "" {
			return fmt.Errorf("fixture must have either away_teamsheet or away_team set")
		}
	}
	return nil
}
