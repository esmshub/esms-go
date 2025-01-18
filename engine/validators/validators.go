package validators

import (
	"fmt"
	"slices"
	"strings"

	"github.com/esmshub/esms-go/engine/models"
	"github.com/esmshub/esms-go/engine/types"
	"go.uber.org/zap"
)

type Validator interface {
	Validate(*models.TeamConfig) error
}

type TeamConfigValidator struct {
	supportedTactics   []string
	supportedPositions []string
	MinSubs            int
	MaxSubs            int
	MinDF              int
	MaxDF              int
	MaxDM              int
	MinMF              int
	MaxMF              int
	MaxAM              int
	MinFW              int
	MaxFW              int
}

func IsValidPosition(position string) bool {
	return slices.ContainsFunc(types.ValidPositions, func(pos string) bool {
		return strings.EqualFold(position, pos)
	})
}

func IsValidTactic(tactic string) bool {
	return slices.ContainsFunc(types.ValidTactics, func(t string) bool {
		return strings.EqualFold(tactic, t)
	})
}

func (v *TeamConfigValidator) validateStartingLineup(players []*models.PlayerPosition) error {
	if len(players) != 11 {
		return fmt.Errorf("expected 11 first team players, got %d", len(players))
	}

	for idx, p := range players {
		zap.L().Debug("validate player position", zap.String("pos", p.Position))
		if !IsValidPosition(p.Position) {
			return fmt.Errorf("player %s has an illegal position (%s)", p.Name, p.Position)
		} else if idx == 0 && p.Position != types.POSITION_GK {
			return fmt.Errorf("player 1 must be a GK")
		} else if p.Position == types.POSITION_GK && idx > 0 {
			return fmt.Errorf("player %s cannot be a GK", p.Name)
		} else if p.Stats.IsInjured {
			return fmt.Errorf("player %s is injured", p.Name)
		} else if p.Stats.IsSuspended {
			return fmt.Errorf("player %s is suspended", p.Name)
		}
	}

	return nil
}

func (v *TeamConfigValidator) validateSubs(subs []*models.PlayerPosition) error {
	count := len(subs)
	if count < v.MinSubs {
		return fmt.Errorf("expected at least %d subs, found %d", v.MinSubs, count)
	} else if count > v.MaxSubs {
		return fmt.Errorf("expected at most %d subs, found %d", v.MaxSubs, count)
	}

	return nil
}

func (v *TeamConfigValidator) validateFormation(tc *models.TeamConfig) error {
	positionCount := make(map[string]int)
	for _, p := range tc.Lineup {
		positionCount[p.Position]++
	}
	defMids := positionCount[types.POSITION_DM]
	mids := positionCount[types.POSITION_MF]
	attMids := positionCount[types.POSITION_AM]

	if positionCount[types.POSITION_DF] < v.MinDF {
		return fmt.Errorf("cannot field less than %d %ss", v.MinDF, types.POSITION_DF)
	} else if positionCount[types.POSITION_DF] > v.MaxDF {
		return fmt.Errorf("cannot field more than %d %ss", v.MaxDF, types.POSITION_DF)
	} else if positionCount[types.POSITION_DM] > v.MaxDM {
		return fmt.Errorf("cannot field more than %d %ss", v.MaxDM, types.POSITION_DM)
	} else if positionCount[types.POSITION_MF] > v.MaxMF {
		return fmt.Errorf("cannot field more than %d %ss", v.MaxMF, types.POSITION_MF)
	} else if positionCount[types.POSITION_AM] > v.MaxAM {
		return fmt.Errorf("cannot field more than %d %ss", v.MaxAM, types.POSITION_AM)
	} else if (mids + defMids + attMids) < v.MinMF {
		return fmt.Errorf("cannot field less than %d midfielders", v.MinMF)
	} else if (mids + defMids + attMids) > v.MaxMF {
		return fmt.Errorf("cannot field more than %d midfielders", v.MaxMF)
	} else if positionCount[types.POSITION_FW] < v.MinFW {
		return fmt.Errorf("cannot field less than %d %ss", v.MinFW, types.POSITION_FW)
	} else if positionCount[types.POSITION_FW] > v.MaxFW {
		return fmt.Errorf("cannot field more than %d %ss", v.MaxFW, types.POSITION_FW)
	}

	formationStr := fmt.Sprintf("%d-%d-%d-%d-%d %s", positionCount[types.POSITION_DF], defMids, mids, attMids, positionCount[types.POSITION_FW], tc.Tactic)
	tc.Formation = strings.ReplaceAll(formationStr, "-0", "")

	return nil
}

func (v *TeamConfigValidator) Validate(tc *models.TeamConfig) error {
	if !IsValidTactic(tc.Tactic) {
		return fmt.Errorf("tactic '%s' is not a supported, expected one of %#v", tc.Tactic, v.supportedTactics)
	}

	if err := v.validateStartingLineup(tc.Lineup); err != nil {
		return err
	}

	if err := v.validateSubs(tc.Subs); err != nil {
		return err
	}

	return nil
}

func NewTeamConfigValidator(config map[string]any) Validator {
	zap.L().Info("config", zap.Any("config", config))
	return &TeamConfigValidator{
		supportedTactics:   types.ValidTactics,
		supportedPositions: types.ValidPositions,
		MinSubs:            config["min_subs"].(int),
		MaxSubs:            config["max_subs"].(int),
		MinDF:              config["min_df"].(int),
		MaxDF:              config["max_df"].(int),
		MaxDM:              config["max_dm"].(int),
		MinMF:              config["min_mf"].(int),
		MaxMF:              config["max_mf"].(int),
		MaxAM:              config["max_am"].(int),
		MinFW:              config["min_fw"].(int),
		MaxFW:              config["max_fw"].(int),
	}
}
