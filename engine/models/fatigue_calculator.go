package models

import (
	"fmt"

	"github.com/esmshub/esms-go/engine/events"
	"github.com/esmshub/esms-go/engine/internal/formulas"
	"github.com/esmshub/esms-go/engine/pkg/utils"
	"go.uber.org/zap"
)

type MatchFatigueCalculator struct {
}

func (m *MatchFatigueCalculator) Update(event string, subject Subject) {
	match, ok := subject.(*Match)
	if !ok {
		zap.L().DPanic("invalid subject type", zap.String("type", fmt.Sprintf("%T", subject)))
	}

	if event == events.MinuteElapsedEventName {
		utils.EachFunc(match.GetTeams(), m.Visit)
	}
}

func (m *MatchFatigueCalculator) Visit(team *MatchTeam) {
	for _, p := range team.GetActive() {
		// oldCond := p.Ability.Condition
		p.SetCondition(p.GetCondition() * formulas.ConditionDecayPerMinute)
		// p.SetCondition(math.Pow(formulas.ConditionDecayPerMinute, float64(p.GetStats().MinutesPlayed)) * p.BaseAbility.Condition)
		// zap.L().Info("player condition", zap.String("name", p.Name), zap.Float64("old", oldCond), zap.Float64("new", p.Ability.Condition))
	}
}
