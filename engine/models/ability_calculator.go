package models

import (
	"fmt"
	"strings"

	"github.com/esmshub/esms-go/engine/events"
	"github.com/esmshub/esms-go/engine/internal/formulas"
	"github.com/esmshub/esms-go/engine/pkg/utils"
	"github.com/esmshub/esms-go/engine/types"
	"go.uber.org/zap"
)

type AbilityCalculator struct {
	tactics           *TacticsMatrix
	defaultAggression int
}

func (m *AbilityCalculator) Update(event string, subject Subject) {
	match, ok := subject.(*Match)
	if !ok {
		zap.L().DPanic("invalid subject type", zap.String("type", fmt.Sprintf("%T", subject)))
	}

	if event == events.MinuteElapsedEventName {
		utils.EachFunc(match.GetTeams(), m.VisitTeam)
	}
}

func (m *AbilityCalculator) VisitTeam(team *MatchTeam) {
	teamAbs := &PlayerAbilities{}
	teamAbs.Goalkeeping = 0
	teamAbs.Tackling = 0
	teamAbs.Passing = 0
	teamAbs.Shooting = 0
	teamAbs.Aggression = 0

	for _, p := range team.GetLineup() {
		playerAbs := &PlayerAbilities{}
		// apply tactical bonuses
		if p.IsActive {
			if p.GetPosition() == types.PositionGK {
				teamAbs.Goalkeeping = p.GetBaseAbility().Goalkeeping
			} else if m.tactics != nil {
				matrix := (*m.tactics)[fmt.Sprintf("%s_%s", team.GetTactic(), p.GetPosition())]
				if matrix == nil || len(matrix) != 3 {
					zap.L().Debug("invalid tactic matrix", zap.String("tactic", team.GetTactic()), zap.String("position", p.GetPosition()))
					matrix = []float64{1, 1, 1}
				}
				baseAbs := p.GetBaseAbility()
				condition := p.GetCondition()
				playerAbs.Tackling = int((matrix[0] * float64(baseAbs.Tackling)) * condition)
				playerAbs.Passing = int((matrix[1] * float64(baseAbs.Passing)) * condition)
				playerAbs.Shooting = int((matrix[2] * float64(baseAbs.Shooting)) * condition)
			}
			// update cumulative stats
			teamAbs.Tackling += playerAbs.Tackling
			teamAbs.Passing += playerAbs.Passing
			teamAbs.Shooting += playerAbs.Shooting
			teamAbs.Aggression += playerAbs.Aggression
		} else {
			playerAbs.Goalkeeping = 0
			playerAbs.Tackling = 0
			playerAbs.Passing = 0
			playerAbs.Shooting = 0
		}

		p.SetMatchAbility(playerAbs)
	}

	agg := m.defaultAggression
	if aggCond := utils.FindFunc(team.GetConditionals(), func(c *Conditional) bool {
		return strings.EqualFold(c.Action, types.AggressionAction)
	}); aggCond != nil {
		agg = aggCond.Values[0].(int)
	}
	teamAbs.Aggression = int(formulas.CalcAggression(float64(agg), float64(teamAbs.Aggression)))
	team.SetAbility(teamAbs)
}

func NewAbilityCalculator(tactics *TacticsMatrix, defaultAggression int) *AbilityCalculator {
	return &AbilityCalculator{
		tactics:           tactics,
		defaultAggression: defaultAggression,
	}
}
