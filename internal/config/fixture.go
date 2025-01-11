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

	zap.L().Info("reading roster file", zap.String("file", rosterFile))
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

	zap.L().Info("reading teamsheet file", zap.String("file", teamsheetFile))
	var config *models.TeamConfig
	if filepath.Ext(teamsheetFile) == ".txt" {
		// load legacy teamsheet
		config, err = LoadLegacyTeamsheet(teamsheetFile, findPlayer)
	} else {
		// TODO: support newer formats
		return nil, fmt.Errorf("unsupported teamsheet file format: %s", teamsheetFile)
	}

	if err != nil {
		return nil, err
	}

	return config, nil
}
