package config

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
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
	GetStringMap(string) map[string]any
	Set(string, any)
	MergeConfigMap(map[string]any) error
	AllSettings() map[string]any
}

var LeagueConfigSupportedFormats = []string{".yaml", ".json", ".dat"}

const DefaultLeagueConfigFileName = "league"
const DefaultMatchReportOutputFileExt = ".txt"

var LeagueConfig Config = viper.New()

func init() {
	exePath, err := os.Executable()
	if err != nil {
		zap.L().Error("unable to get current working directory", zap.Error(err))
	}

	conf := LeagueConfig.(*viper.Viper)
	conf.SetDefault("name", "ESMS League")
	conf.SetDefault("paths.roster_dir", exePath)
	conf.SetDefault("paths.fixtureset_dir", exePath)
	conf.SetDefault("paths.teamsheet_dir", exePath)
	conf.SetDefault("paths.output_dir", exePath)
	conf.SetDefault("match.home_bonus", 100)
	conf.SetDefault("match.extra_time", false)
	conf.SetDefault("match.min_subs", 3)
	conf.SetDefault("match.max_subs", 7)
	conf.SetDefault("match.min_df", 3)
	conf.SetDefault("match.max_df", 5)
	conf.SetDefault("match.max_dm", 3)
	conf.SetDefault("match.min_mf", 1)
	conf.SetDefault("match.max_mf", 6)
	conf.SetDefault("match.max_am", 3)
	conf.SetDefault("match.min_fw", 0)
	conf.SetDefault("match.max_fw", 4)
	conf.SetDefault("match.commentary_file", path.Join(GetConfigDir(), "language.dat"))
	// conf.SetDefault("bonuses.goal", 30)
	// conf.SetDefault("bonuses.assist", 21)
	// conf.SetDefault("bonuses.victory", 30)
	// conf.SetDefault("bonuses.clean_sheet", 20)
	// conf.SetDefault("bonuses.key_tackle", 15)
	// conf.SetDefault("bonuses.key_pass", 12)
	// conf.SetDefault("bonuses.shot_on_target", 8)
	// conf.SetDefault("bonuses.shot_off_target", 0)
	// conf.SetDefault("bonuses.save", 10)
	// conf.SetDefault("bonuses.own_goal", -10)
	// conf.SetDefault("bonuses.defeat", -30)
	// conf.SetDefault("bonuses.conceded", -8)
	// conf.SetDefault("bonuses.cautioned", -3)
	// conf.SetDefault("bonuses.sent_off", -10)
	// legacy support
	conf.RegisterAlias("games", "name")
	conf.RegisterAlias("abbreviations", "teams")
	conf.RegisterAlias("home_bonus", "match.home_bonus")
	conf.RegisterAlias("extra_time", "match.extra_time")
	conf.RegisterAlias("min_subs", "match.min_subs")
	conf.RegisterAlias("bench_size", "match.max_subs")
	conf.RegisterAlias("min_df", "match.min_df")
	conf.RegisterAlias("max_df", "match.max_df")
	conf.RegisterAlias("max_dm", "match.max_dm")
	conf.RegisterAlias("min_mf", "match.min_mf")
	conf.RegisterAlias("max_mf", "match.max_mf")
	conf.RegisterAlias("max_am", "match.max_am")
	conf.RegisterAlias("min_fw", "match.min_fw")
	conf.RegisterAlias("max_fw", "match.max_fw")
	conf.RegisterAlias("bonuses", "abilities")
	// conf.RegisterAlias("bonuses.save", "abilities.ab_sav")
	// conf.RegisterAlias("bonuses.goal", "abilities.ab_goal")
	// conf.RegisterAlias("bonuses.assist", "abilities.ab_assist")
	// conf.RegisterAlias("bonuses.victory", "abilities.ab_victory_random")
	// conf.RegisterAlias("bonuses.clean_sheet", "abilities.ab_clean_sheet")
	// conf.RegisterAlias("bonuses.key_tackle", "abilities.ab_ktk")
	// conf.RegisterAlias("bonuses.key_pass", "abilities.ab_kps")
	// conf.RegisterAlias("bonuses.shot_on_target", "abilities.ab_sht_on")
	// conf.RegisterAlias("bonuses.shot_off_target", "abilities.ab_sht_off")
	// conf.RegisterAlias("bonuses.own_goal", "abilities.ab_og")
	// conf.RegisterAlias("bonuses.defeat", "abilities.ab_defeat_random")
	// conf.RegisterAlias("bonuses.conceded", "abilities.ab_concede")
	// conf.RegisterAlias("bonuses.cautioned", "abilities.ab_yellow")
	// conf.RegisterAlias("bonuses.sent_off", "abilities.ab_red")
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
	zap.L().Debug("Root paths", zap.Strings("paths", rootPaths))
	for _, ext := range LeagueConfigSupportedFormats {
		for _, dir := range rootPaths {
			configFilePath := filepath.Join(dir, fmt.Sprintf("%s%s", DefaultLeagueConfigFileName, ext))
			zap.L().Debug("Checking for config file", zap.String("path", configFilePath))
			if utils.FileExists(configFilePath) {
				return LoadLeagueConfig(configFilePath)
			}
		}
	}

	return fmt.Errorf("no league config file found in %s", GetConfigDir())
}

func LoadLeagueConfig(filePath string) error {
	var err error
	fileExt := filepath.Ext(filePath)
	if fileExt == "" {
		return errors.New("config file extension is missing")
	} else if fileExt == ".dat" {
		// treat DAT file as properties format
		zap.L().Info("Reading league config", zap.String("path", filePath))
		zap.L().Warn("Legacy DAT file detected")
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
					if intValue, err := strconv.Atoi(value); err == nil {
						LeagueConfig.Set(key, intValue)
					} else {
						LeagueConfig.Set(key, value)
					}
				} else {
					// Handle invalid key-value pair format
					return fmt.Errorf("invalid row format: %s", line)
				}
			}

			return nil
		})
	} else {
		zap.L().Info("Reading league config", zap.String("path", filePath))
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
