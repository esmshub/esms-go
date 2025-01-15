/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/esmshub/esms-go/engine"
	"github.com/esmshub/esms-go/engine/models"
	"github.com/esmshub/esms-go/internal/config"
	"github.com/esmshub/esms-go/internal/store"
	"github.com/esmshub/esms-go/pkg/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	fixtureSetFilePath string
	rngSeed            uint64
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "sim",
	Short: "Run a simulation",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configFilePath, err := cmd.Flags().GetString("config-file")
		if err != nil {
			return err
		}

		if configFilePath != "" {
			err = config.LoadLeagueConfig(configFilePath)
		} else {
			zap.L().Info("loading nearest league config")
			err = config.LoadNearestLeagueConfig()
		}

		if err != nil {
			zap.L().Warn("unable to load config", zap.Error(err))
			zap.L().Warn("using default config")
		}

		fixtureSet := utils.Must(config.LoadFixtureset[config.Fixtureset](fixtureSetFilePath))
		if fixtureSet.OverrideSettings != nil {
			zap.L().Debug("applying override settings", zap.Any("settings", fixtureSet.OverrideSettings))
			err := config.LeagueConfig.MergeConfigMap(fixtureSet.OverrideSettings)
			if err != nil {
				zap.L().Warn("unable to apply override settings", zap.Error(err))
			}
		}

		opts := &engine.Options{
			HomeBonus:       config.LeagueConfig.GetInt("match.home_bonus"),
			EnableExtraTime: config.LeagueConfig.GetBool("match.extra_time"),
			RngSeed:         rngSeed,
		}
		zap.L().Info("config", zap.Any("config", config.LeagueConfig.AllSettings()))
		for _, f := range fixtureSet.Fixtures {
			// resolve config paths
			if f.HomeTeamCode != "" {
				f.HomeTeamsheet = fmt.Sprintf("%ssht%s", filepath.Join(config.LeagueConfig.GetString("paths.teamsheet_dir"), f.HomeTeamCode), config.DefaultTeamsheetFileExt)
				f.HomeRoster = fmt.Sprintf("%s%s", filepath.Join(config.LeagueConfig.GetString("paths.roster_dir"), f.HomeTeamCode), config.DefaultRosterFileExt)
			}
			if f.AwayTeamCode != "" {
				f.AwayTeamsheet = fmt.Sprintf("%ssht%s", filepath.Join(config.LeagueConfig.GetString("paths.teamsheet_dir"), f.AwayTeamCode), config.DefaultTeamsheetFileExt)
				f.AwayRoster = fmt.Sprintf("%s%s", filepath.Join(config.LeagueConfig.GetString("paths.roster_dir"), f.AwayTeamCode), config.DefaultRosterFileExt)
			}
			fixturesetDir := filepath.Dir(fixtureSetFilePath)
			if !filepath.IsAbs(f.HomeTeamsheet) {
				f.HomeTeamsheet = filepath.Join(fixturesetDir, f.HomeTeamsheet)
			}
			if !filepath.IsAbs(f.HomeRoster) {
				f.HomeRoster = filepath.Join(fixturesetDir, f.HomeRoster)
			}
			if !filepath.IsAbs(f.AwayTeamsheet) {
				f.AwayTeamsheet = filepath.Join(fixturesetDir, f.AwayTeamsheet)
			}
			if !filepath.IsAbs(f.AwayRoster) {
				f.AwayRoster = filepath.Join(fixturesetDir, f.AwayRoster)
			}
			// load teams / rosters
			homeTeam := utils.Must(config.LoadTeamConfig(f.HomeTeamsheet, f.HomeRoster))
			homeTeam.Name = config.LeagueConfig.Get("teams").(map[string]any)[homeTeam.Code].(string)
			homeTeam.ManagerName = config.LeagueConfig.Get("managers").(map[string]any)[homeTeam.Code].(string)
			homeTeam.StadiumName = config.LeagueConfig.Get("stadiums").(map[string]any)[homeTeam.Code].(string)
			homeTeam.StadiumCapacity = utils.Must(strconv.Atoi(config.LeagueConfig.Get("stadium_capacity").(map[string]any)[homeTeam.Code].(string)))
			awayTeam := utils.Must(config.LoadTeamConfig(f.AwayTeamsheet, f.AwayRoster))
			awayTeam.Name = config.LeagueConfig.Get("teams").(map[string]any)[awayTeam.Code].(string)
			awayTeam.ManagerName = config.LeagueConfig.Get("managers").(map[string]any)[awayTeam.Code].(string)
			awayTeam.StadiumName = config.LeagueConfig.Get("stadiums").(map[string]any)[awayTeam.Code].(string)
			awayTeam.StadiumCapacity = utils.Must(strconv.Atoi(config.LeagueConfig.Get("stadium_capacity").(map[string]any)[awayTeam.Code].(string)))
			match := &models.Match{
				HomeTeam: homeTeam,
				AwayTeam: awayTeam,
			}

			zap.L().Info("running fixture", zap.Any("fixture", f))
			result := engine.Run(match, opts)
			fileStore := store.MatchResultFileStore{}
			err := fileStore.Save(result, store.MatchResultFileStoreOptions{
				LeagueName: config.LeagueConfig.GetString("name"),
				OutputFile: filepath.Join(config.LeagueConfig.GetString("paths.output_dir"), fmt.Sprintf("%s_%s%s", homeTeam.Code, awayTeam.Code, config.DefaultMatchReportOutputFileExt)),
				RngSeed:    result.RngSeed,
			})
			if err != nil {
				zap.L().Error("unable to save match result", zap.Error(err))
			}
			fmt.Println("------------------------------")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&fixtureSetFilePath, "fixture-set", "f", "", "Path to fixture set file")
	runCmd.Flags().Uint64VarP(&rngSeed, "rng-seed", "s", 0, "Seed for random number generator")

	runCmd.MarkFlagRequired("fixture-set")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
