package config

import "github.com/esmshub/esms-go/pkg/utils"

// MergeWithDefaults merges a given map into the current Viper state (preserving defaults).
func MergeWithDefaults(override map[string]interface{}) error {
	base := LeagueConfig.AllSettings() // includes SetDefault values
	merged := utils.DeepMerge(base, override)
	return LeagueConfig.MergeConfigMap(merged)
}
