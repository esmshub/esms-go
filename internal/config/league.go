package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/esmshub/esms-go/pkg/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config interface {
	GetInt(string) int
	GetString(string) string
	GetBool(string) bool
	Get(string) any
	Set(string, any)
	MergeConfigMap(map[string]any) error
	AllSettings() map[string]any
}

var LeagueConfigSupportedFormats = []string{".yaml", ".json", ".dat"}

const DefaultLeagueConfigFileName = "league"

var LeagueConfig Config = viper.New()

func init() {
	conf := LeagueConfig.(*viper.Viper)
	conf.SetDefault("home_bonus", 30)
	conf.SetDefault("extra_time", false)
	conf.SetDefault("min_subs", 3)
	conf.SetDefault("bench_size", 7)
	conf.SetDefault("min_df", 3)
	conf.SetDefault("max_df", 5)
	conf.SetDefault("max_dm", 3)
	conf.SetDefault("min_mf", 1)
	conf.SetDefault("max_mf", 6)
	conf.SetDefault("max_am", 3)
	conf.SetDefault("min_fw", 0)
	conf.SetDefault("max_fw", 4)
}

func LoadNearestLeagueConfig() error {
	// Get the path to the executable
	exePath, err := os.Executable()
	if err != nil {
		zap.L().Warn("unable to get the executable path", zap.Error(err))
	}

	rootPaths := []string{exePath}
	if configDir := GetConfigDir(); configDir != "" {
		rootPaths = append(rootPaths, configDir)
	}
	zap.L().Info("Root paths", zap.Strings("paths", rootPaths))
	for _, ext := range LeagueConfigSupportedFormats {
		for _, dir := range rootPaths {
			configFilePath := filepath.Join(dir, fmt.Sprintf("%s%s", DefaultLeagueConfigFileName, ext))
			zap.L().Info("Checking for config file", zap.String("path", configFilePath))
			if utils.FileExists(configFilePath) {
				return LoadLeagueConfig(configFilePath)
			}
		}
	}

	return fmt.Errorf("no league config file found in %s", GetConfigDir())
}

func LoadLeagueConfig(filePath string) error {
	var err error
	fmt.Printf("Reading config from %s\n", filePath)
	fileExt := filepath.Ext(filePath)
	if fileExt == "" {
		return errors.New("config file extension is missing")
	} else if fileExt == ".dat" {
		// treat DAT file as properties format
		zap.L().Warn("Legacy DAT file detected")
		zap.L().Info("Reading league config...")
		prefix := ""
		_, err = utils.ReadFile(filePath, func(line string, row int) error {
			if strings.HasSuffix(line, ":") {
				prefix = fmt.Sprintf("%s.", strings.ToLower(strings.Trim(line, ":")))
			} else if len(strings.TrimSpace(line)) == 0 {
				prefix = ""
			} else {
				parts := strings.Split(line, "=")
				if len(parts) == 2 {
					key := fmt.Sprintf("%s%s", prefix, strings.TrimSpace(parts[0]))
					value := strings.TrimSpace(parts[1])
					zap.L().Debug("Setting parsed", zap.String(key, value))
					LeagueConfig.Set(key, value)
				} else {
					// Handle invalid key-value pair format
					return fmt.Errorf("invalid row format: %s", line)
				}
			}

			return nil
		})
	} else {
		zap.L().Info("Reading league config...")
		LeagueConfig.(*viper.Viper).SetConfigFile(filePath)
		err = LeagueConfig.(*viper.Viper).ReadInConfig()
	}

	if err != nil {
		return err
	}

	// Unmarshal the configuration into the Config struct
	// err = config.Unmarshal(&config)
	// if err != nil {
	// 	return config, err
	// }

	return nil
}
