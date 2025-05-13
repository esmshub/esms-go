package config

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/esmshub/esms-go/engine/models"
	"go.uber.org/zap"
)

const DefaultTeamsheetFileExt = ".txt"
const DefaultRosterFileExt = ".txt"

type Fixture struct {
	HomeTeamCode  string `mapstructure:"home_team" json:"home_team" yaml:"home_team"`
	AwayTeamCode  string `mapstructure:"away_team" json:"away_team" yaml:"away_team"`
	HomeTeamsheet string `mapstructure:"home_teamsheet" json:"home_teamsheet" yaml:"home_teamsheet"`
	HomeRoster    string `mapstructure:"home_roster" json:"home_roster" yaml:"home_roster"`
	AwayTeamsheet string `mapstructure:"away_teamsheet" json:"away_teamsheet" yaml:"away_teamsheet"`
	AwayRoster    string `mapstructure:"away_roster" json:"away_roster" yaml:"away_roster"`
}

func LoadTeamConfig(teamsheetFile, rosterFile string) (*models.TeamConfig, error) {
	// load roster
	var roster []*models.Player
	var err error

	zap.L().Info("reading roster file", zap.String("path", rosterFile))
	if filepath.Ext(rosterFile) == ".txt" {
		// load legacy teamsheet
		roster, err = LoadRoster(rosterFile)
	} else {
		// TODO: support newer formats
		return nil, fmt.Errorf("unsupported roster file format: %s", rosterFile)
	}

	if err != nil {
		return nil, err
	}

	findPlayer := func(name string) *models.Player {
		i := slices.IndexFunc(roster, func(p *models.Player) bool {
			return strings.EqualFold(p.Name, name)
		})
		if i == -1 {
			return nil
		}

		return roster[i]
	}

	zap.L().Info("reading teamsheet", zap.String("path", teamsheetFile))
	var config *models.TeamConfig
	if filepath.Ext(teamsheetFile) == ".txt" {
		// load legacy teamsheet
		config, err = LoadLegacyTeamsheet(teamsheetFile, findPlayer)
		if err == nil && config.Code == "" {
			config.Code = strings.TrimSuffix(filepath.Base(teamsheetFile), DefaultTeamsheetFileExt)
		}
	} else {
		// TODO: support newer formats
		return nil, fmt.Errorf("unsupported teamsheet file format: %s", teamsheetFile)
	}

	if err != nil {
		return nil, err
	}

	teamsMap := LeagueConfig.GetStringMap("teams")
	teamName, nameOk := teamsMap[config.Code].(string)
	if !nameOk {
		return nil, fmt.Errorf("teams.%s config not set", config.Code)
	}
	config.Name = teamName
	managersMap := LeagueConfig.GetStringMap("managers")
	managerName, managerOk := managersMap[config.Code].(string)
	if !managerOk {
		return nil, fmt.Errorf("managers.%s config not set", config.Code)
	}
	config.ManagerName = managerName
	stadiumsMap := LeagueConfig.GetStringMap("stadiums")
	stadiumName, stadiumOk := stadiumsMap[config.Code].(string)
	if !stadiumOk {
		return nil, fmt.Errorf("stadiums.%s config not set", config.Code)
	}
	config.StadiumName = stadiumName
	stadiumCapMap := LeagueConfig.GetStringMap("capacities")
	stadiumCap, stadiumCapOk := stadiumCapMap[config.Code].(int)
	if !stadiumCapOk {
		return nil, fmt.Errorf("capacities.%s config not set", config.Code)
	}
	config.StadiumCapacity = stadiumCap

	return config, nil
}
