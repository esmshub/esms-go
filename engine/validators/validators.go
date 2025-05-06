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
	Validate(*models.MatchTeam) error
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

func (v *TeamConfigValidator) validateStartingLineup(players []*models.MatchPlayer) error {
	if len(players) != 11 {
		return fmt.Errorf("expected 11 first team players, got %d", len(players))
	}

	for idx, p := range players {
		zap.L().Debug("validate player position", zap.String("pos", p.GetPosition()))
		if !IsValidPosition(p.GetPosition()) {
			return fmt.Errorf("player %s has an illegal position (%s)", p.GetName(), p.GetPosition())
		} else if idx == 0 && p.GetPosition() != types.PositionGK {
			return fmt.Errorf("player 1 must be a GK")
		} else if p.GetPosition() == types.PositionGK && idx > 0 {
			return fmt.Errorf("player %s cannot be a GK", p.GetName())
		} else if p.GetIsInjured() {
			return fmt.Errorf("player %s is injured", p.GetName())
		} else if p.GetIsSuspended() {
			return fmt.Errorf("player %s is suspended", p.GetName())
		}
	}

	return nil
}

func (v *TeamConfigValidator) validateSubs(subs []*models.MatchPlayer) error {
	count := len(subs)
	if count < v.MinSubs {
		return fmt.Errorf("expected at least %d subs, found %d", v.MinSubs, count)
	} else if count > v.MaxSubs {
		return fmt.Errorf("expected at most %d subs, found %d", v.MaxSubs, count)
	}

	return nil
}

func (v *TeamConfigValidator) validateFormation(team *models.MatchTeam) error {
	positionCount := make(map[string]int)
	for _, p := range team.GetStarters() {
		positionCount[p.GetPosition()]++
	}
	defMids := positionCount[types.PositionDM]
	mids := positionCount[types.PositionMF]
	attMids := positionCount[types.PositionAM]

	if positionCount[types.PositionDF] < v.MinDF {
		return fmt.Errorf("cannot field less than %d %ss", v.MinDF, types.PositionDF)
	} else if positionCount[types.PositionDF] > v.MaxDF {
		return fmt.Errorf("cannot field more than %d %ss", v.MaxDF, types.PositionDF)
	} else if positionCount[types.PositionDM] > v.MaxDM {
		return fmt.Errorf("cannot field more than %d %ss", v.MaxDM, types.PositionDM)
	} else if positionCount[types.PositionMF] > v.MaxMF {
		return fmt.Errorf("cannot field more than %d %ss", v.MaxMF, types.PositionMF)
	} else if positionCount[types.PositionAM] > v.MaxAM {
		return fmt.Errorf("cannot field more than %d %ss", v.MaxAM, types.PositionAM)
	} else if (mids + defMids + attMids) < v.MinMF {
		return fmt.Errorf("cannot field less than %d midfielders", v.MinMF)
	} else if (mids + defMids + attMids) > v.MaxMF {
		return fmt.Errorf("cannot field more than %d midfielders", v.MaxMF)
	} else if positionCount[types.PositionFW] < v.MinFW {
		return fmt.Errorf("cannot field less than %d %ss", v.MinFW, types.PositionFW)
	} else if positionCount[types.PositionFW] > v.MaxFW {
		return fmt.Errorf("cannot field more than %d %ss", v.MaxFW, types.PositionFW)
	}

	return nil
}

func (v *TeamConfigValidator) Validate(team *models.MatchTeam) error {
	if !IsValidTactic(team.GetTactic()) {
		return fmt.Errorf("tactic '%s' is not a supported, expected one of %#v", team.GetTactic(), v.supportedTactics)
	}

	if err := v.validateStartingLineup(team.GetStarters()); err != nil {
		return err
	}

	if err := v.validateSubs(team.GetSubs()); err != nil {
		return err
	}

	if err := v.validateFormation(team); err != nil {
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
