package config

import (
	"os"

	"github.com/esmshub/esms-go/pkg/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const ConfigFolderName = ".esms"

var Cli Config = viper.New()

func init() {
	exePath, err := os.Executable()
	if err != nil {
		zap.L().Error("unable to get current working directory", zap.Error(err))
	}
	conf := Cli.(*viper.Viper)
	conf.SetDefault("paths.roster_dir", exePath)
	conf.SetDefault("paths.fixtureset_dir", exePath)
	conf.SetDefault("paths.teamsheet_dir", exePath)
	conf.SetDefault("paths.config_dir", exePath)
}

func GetConfigDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		zap.L().Warn("unable to get the current directory", zap.Error(err))
		return "" // default to current folder
	}

	path, err := utils.FindAncestorDir(cwd, ConfigFolderName, true)
	if err != nil {
		zap.L().Warn("an error occurred whilst traversing dir tree", zap.Error(err))
		return ""
	} else if path == "" {
		// check user profile / home dir
		homeDir, err := os.UserHomeDir()
		if err != nil {
			zap.L().Warn("unable to get the home directory", zap.Error(err))
			return ""
		}

		path, err = utils.FindAncestorDir(homeDir, ConfigFolderName, false)
		if err != nil {
			zap.L().Warn("an error occurred reading home directory", zap.Error(err))
		} else if path == "" {
			return ""
		}

		return path
	}

	return path
}
