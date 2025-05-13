package config

import (
	"os"

	"github.com/esmshub/esms-go/pkg/utils"
	"go.uber.org/zap"
)

const ConfigFolderName = ".esms"

func GetConfigDir() string {
	exePath, err := os.Executable()
	if err != nil {
		zap.L().Warn("unable to get the executable path", zap.Error(err))
		exePath = "." // default to current folder
	}

	path, err := utils.FindAncestorDir(exePath, ConfigFolderName, true)
	if err != nil {
		zap.L().Warn("an error occurred whilst traversing dir tree", zap.Error(err))
		return "."
	} else if path == "" {
		// check user profile / home dir
		homeDir, err := os.UserHomeDir()
		if err != nil {
			zap.L().Warn("unable to get the home directory", zap.Error(err))
			return exePath
		}

		path, err = utils.FindAncestorDir(homeDir, ConfigFolderName, false)
		if err != nil {
			zap.L().Warn("an error occurred reading home directory", zap.Error(err))
		} else if path == "" {
			path = exePath
		}

		return path
	}

	return path
}
