/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"maps"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/esmshub/esms-go/engine"
	"github.com/esmshub/esms-go/engine/commentary"
	"github.com/esmshub/esms-go/engine/models"
	"github.com/esmshub/esms-go/internal/config"
	"github.com/esmshub/esms-go/internal/formatters"
	"github.com/esmshub/esms-go/internal/store"
	"github.com/esmshub/esms-go/pkg/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	fixtureSetFilePath string
	tacticsFilePath    string
	rngSeed            uint64
	contentStyle       = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFF"))
		// Background(lipgloss.Color("#000000")).

	headerStyle   = contentStyle.Foreground(lipgloss.Color("#0f0"))
	teamNameStyle = contentStyle.Foreground(lipgloss.Color("#0FF"))
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "sim",
	Short: "Simulate a round of fixtures",
	Long: `Simulate each match from the provided fixture set and generate a match report.

A valid fixture set file must contain a name and at least one fixture. Each fixture must 
define either the teamsheet filename or the team code of the home and away side. 

When specifying the teamsheet filename, paths should be relative to the fixture set file. 
When specifying the team code, the path and format of the teamsheet file will be assumed
based on the default config i.e. <paths.teamsheet_dir>/<team_code>sht.txt

Fixture set files can be in either JSON or YAML format e.g.

fixtures.yml
----------------------------
name: Round 1
fixtures:
  - home_teamsheet: ransht.txt
    away_teamsheet: celsht.txt

fixtures.json
----------------------------
{
  "name": "Round 1",
  "fixtures": [
    {
      "home_team": "ran",
      "away_team": "cel"
    }
  ]
}
----------------------------

If a tactics matrix file is provided, appropriate bonuses will be applied otherwise, no 
tactical bonuses are used. If no '-c/--config-file' flag is provided, the nearest valid 
config file will be loaded.

Configuration settings can overridden at fixture set level by declaring an 'override_settings' block
in the fixture set file. This block should be a map of config settings, for example:

config.yml
---------------------------------
paths:
  teamsheet_dir: /league/teamsheets
match:
  extra_time: false


fixtures.yml
---------------------------------
name: Cup Quarter Final - 2nd Leg
fixtures:
  - home_teamsheet: ransht.txt
    away_teamsheet: celsht.txt
override_settings:
  paths:
    teamsheet_dir: /cup/teamsheets
  match:
    extra_time: true
---------------------------------

In the above example, teamsheets will be read from /cup/teamsheets and extra time will 
be enabled this fixture set only.
`,
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

		fixtureSet := utils.Must(config.LoadFixtureset(fixtureSetFilePath))
		if fixtureSet.OverrideSettings != nil {
			zap.L().Debug("applying override settings", zap.Any("settings", maps.Keys(fixtureSet.OverrideSettings)))
			err := config.MergeWithDefaults(fixtureSet.OverrideSettings)
			if err != nil {
				zap.L().Warn("unable to apply override settings", zap.Error(err))
			}
		}

		if tacticsFilePath == "" {
			tacticsFilePath = filepath.Join(config.GetConfigDir(), config.DefaultTacticsMatrixFileName)
		}

		zap.L().Info("Reading tactics matrix", zap.String("path", tacticsFilePath))
		tacticsMatrix, err := config.LoadTactics(tacticsFilePath)
		if err != nil {
			zap.L().Warn("unable to load tactics matrix", zap.Error(err))
		}
		zap.L().Info("Reading commentary file", zap.String("path", config.LeagueConfig.GetString("match.commentary_file")))
		commsProvider := commentary.NewLegacyFileCommentaryProvider()
		if err := commsProvider.Load(config.LeagueConfig.GetString("match.commentary_file")); err != nil {
			zap.L().Warn("unable to load commentary provider", zap.Error(err))
		}

		opts := &engine.Options{
			RngSeed:            rngSeed,
			TacticsMatrix:      tacticsMatrix,
			AppConfig:          config.LeagueConfig.GetStringMap("match"),
			CommentaryProvider: commsProvider,
		}

		t := table.New().
			Headers(fmt.Sprintf("%s RESULTS", strings.ToUpper(fixtureSet.Name))).
			// BorderHeader(false).
			Border(lipgloss.HiddenBorder()).
			// Border(lipgloss.RoundedBorder()).
			// BorderStyle(lipgloss.NewStyle().Background(lipgloss.Color("#000"))).
			// BorderColumn(false).
			// BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("240")))
			// BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("240"))).
			StyleFunc(func(row, col int) lipgloss.Style {
				if row == -1 {
					return headerStyle
				} else if col != 1 {
					return teamNameStyle
				}

				return contentStyle
			}).
			Rows([]string{})
		fmtScorerOpts := formatters.FormatScorersOptions{
			RowDelimiter: "\n",
		}
		for fid, f := range fixtureSet.Fixtures {
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
			homeTeam, homeConfErr := config.LoadTeamConfig(f.HomeTeamsheet, f.HomeRoster)
			if homeConfErr != nil {
				zap.L().Panic("unable to load home team config", zap.Error(homeConfErr))
			}
			awayTeam, awayConfErr := config.LoadTeamConfig(f.AwayTeamsheet, f.AwayRoster)
			if awayConfErr != nil {
				zap.L().Panic("unable to load away team config", zap.Error(awayConfErr))
			}
			match := models.NewMatch(homeTeam, awayTeam)

			zap.L().Info(fmt.Sprintf("running fixture %d of %d", fid+1, len(fixtureSet.Fixtures)))
			result, err := engine.Run(match, opts)
			if err != nil {
				zap.L().Panic("unable to run match", zap.Error(err))
			}

			homeGoals := result.HomeTeam.GetStats().Goals
			awayGoals := result.AwayTeam.GetStats().Goals

			homeStyle := teamNameStyle.Padding(0, 0)
			awayStyle := homeStyle
			if len(homeGoals) > len(awayGoals) {
				homeStyle = homeStyle.Foreground(lipgloss.Color("#FF0"))
			} else if len(homeGoals) < len(awayGoals) {
				awayStyle = homeStyle.Foreground(lipgloss.Color("#FF0"))
			}

			t.Row(
				fmt.Sprintf("%s\n%s", homeStyle.Render(strings.ToUpper(result.HomeTeam.GetName())), contentStyle.Padding(0, 0).Render(formatters.FormatScorers(homeGoals, fmtScorerOpts))),
				fmt.Sprintf("%d-%d", len(homeGoals), len(awayGoals)),
				fmt.Sprintf("%s\n%s", awayStyle.Render(strings.ToUpper(result.AwayTeam.GetName())), contentStyle.Padding(0, 0).Render(formatters.FormatScorers(awayGoals, fmtScorerOpts))),
			)
			t.Row()

			// apply bonuses
			bonusConfig := config.LeagueConfig.GetStringMap("bonuses")
			if bonusConfig != nil {
				models.NewMatchBonusCalculator(bonusConfig).Apply(result)
			} else {
				zap.L().Warn("no valid bonus config found")
			}

			comms := []string{}
			if legacyCommentary, ok := opts.CommentaryProvider.(*commentary.LegacyFileCommentaryProvider); ok {
				comms = legacyCommentary.GetCommentary()
				legacyCommentary.Clear()
			}

			fileStore := store.MatchResultFileStore{}
			err = fileStore.Save(result, comms, store.MatchResultFileStoreOptions{
				HeaderText: config.LeagueConfig.GetString("name"),
				OutputFile: filepath.Join(config.LeagueConfig.GetString("paths.output_dir"), fmt.Sprintf("%s_%s%s", homeTeam.Code, awayTeam.Code, config.DefaultMatchReportOutputFileExt)),
				FooterText: fmt.Sprintf("\n%d Produced from %s v%s", result.RngSeed, cmd.Root().Use, cmd.Root().Version),
			})
			if err != nil {
				zap.L().Error("unable to save match result", zap.Error(err))
			}
		}
		fmt.Println(t.Render())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&fixtureSetFilePath, "fixture-set", "f", "", "Path to fixture set file")
	runCmd.Flags().Uint64VarP(&rngSeed, "rng-seed", "s", 0, "Seed for random number generator")
	runCmd.Flags().StringVarP(&tacticsFilePath, "tactics", "t", "", "Path to tactics matrix file")

	runCmd.MarkFlagRequired("fixture-set")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
