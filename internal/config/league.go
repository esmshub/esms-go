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
const DefaultCommentaryFileName = "language.dat"

var LeagueConfig Config = viper.New()

var configKeyAliases = map[string]string{
	"abbrevations": "teams",
	// "abilities":                   "bonuses",
	"abilities.ab_sav":            "bonuses.save",
	"abilities.ab_goal":           "bonuses.goal",
	"abilities.ab_assist":         "bonuses.assist",
	"abilities.ab_victory_random": "bonuses.victory",
	"abilities.ab_clean_sheet":    "bonuses.clean_sheet",
	"abilities.ab_ktk":            "bonuses.key_tackle",
	"abilities.ab_kps":            "bonuses.key_pass",
	"abilities.ab_sht_on":         "bonuses.shot_on_target",
	"abilities.ab_sht_off":        "bonuses.shot_off_target",
	"abilities.ab_og":             "bonuses.own_goal",
	"abilities.ab_defeat_random":  "bonuses.defeat",
	"abilities.ab_concede":        "bonuses.conceded",
	"abilities.ab_yellow":         "bonuses.cautioned",
	"abilities.ab_red":            "bonuses.sent_off",
	"home_bonus":                  "bonuses.home_adv",
	"games":                       "name",
	"extra_time":                  "match.extra_time",
	"num_substitutions":           "match.min_subs",
	"bench_size":                  "match.max_subs",
	"min_df":                      "match.min_df",
	"max_df":                      "match.max_df",
	"max_dm":                      "match.max_dm",
	"min_mf":                      "match.min_mf",
	"max_mf":                      "match.max_mf",
	"max_am":                      "match.max_am",
	"min_fw":                      "match.min_fw",
	"max_fw":                      "match.max_fw",
}

func init() {
	exePath, err := os.Executable()
	if err != nil {
		zap.L().Warn("unable to get the executable path", zap.Error(err))
		exePath = "." // default to current folder
	}

	conf := LeagueConfig.(*viper.Viper)
	conf.SetDefault("name", "ESMS League")
	conf.SetDefault("paths.roster_dir", exePath)
	conf.SetDefault("paths.teamsheet_dir", exePath)
	conf.SetDefault("paths.output_dir", exePath)
	matchConfig := map[string]any{
		"extra_time":      false,
		"min_subs":        3,
		"max_subs":        7,
		"min_df":          3,
		"max_df":          5,
		"max_dm":          3,
		"min_mf":          1,
		"max_mf":          6,
		"max_am":          3,
		"min_fw":          0,
		"max_fw":          4,
		"commentary_file": path.Join(GetConfigDir(), DefaultCommentaryFileName),
	}
	conf.SetDefault("match", matchConfig)
	conf.SetDefault("teams", map[string]string{})
	conf.SetDefault("managers", map[string]string{})
	conf.SetDefault("stadiums", map[string]string{})
	conf.SetDefault("capacities", map[string]int{})
	bonusConfig := map[string]int{
		"home_adv":        100,
		"goal":            30,
		"assist":          21,
		"victory":         30,
		"clean_sheet":     20,
		"key_tackle":      15,
		"key_pass":        12,
		"shot_on_target":  8,
		"shot_off_target": 0,
		"save":            10,
		"own_goal":        -10,
		"defeat":          -30,
		"conceded":        -8,
		"cautioned":       -3,
		"sent_off":        -10,
	}
	conf.SetDefault("bonuses", bonusConfig)
}

func LoadNearestLeagueConfig() error {
	for _, ext := range LeagueConfigSupportedFormats {
		configFilePath := filepath.Join(GetConfigDir(), fmt.Sprintf("%s%s", DefaultLeagueConfigFileName, ext))
		zap.L().Debug("Checking for config file", zap.String("path", configFilePath))
		if utils.FileExists(configFilePath) {
			return LoadLeagueConfig(configFilePath)
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
				rootKey := strings.ToLower(strings.Trim(line, ":"))
				alias, ok := configKeyAliases[strings.ToLower(rootKey)]
				if ok {
					zap.L().Warn("Setting alias", zap.String("old_key", rootKey), zap.String("new_key", alias))
					rootKey = alias
				}
				prefix = fmt.Sprintf("%s.", rootKey)
			} else if len(strings.TrimSpace(line)) == 0 {
				prefix = ""
			} else {
				parts := strings.Split(line, "=")
				if len(parts) == 2 {
					key := fmt.Sprintf("%s%s", prefix, strings.TrimSpace(parts[0]))
					alias, ok := configKeyAliases[strings.ToLower(key)]
					if ok {
						zap.L().Warn("Setting alias", zap.String("old_key", key), zap.String("new_key", alias))
						key = alias
					}
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
