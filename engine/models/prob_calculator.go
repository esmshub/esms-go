package models

import (
	"fmt"

	"github.com/esmshub/esms-go/engine/events"
	"go.uber.org/zap"
)

type ProbabilityCalculator struct{}

func (m *ProbabilityCalculator) Update(event string, subject Subject) {
	match, ok := subject.(*Match)
	if !ok {
		zap.L().DPanic("invalid subject type", zap.String("type", fmt.Sprintf("%T", subject)))
	}

	if event == events.MinuteElapsedEventName {
		match.HomeTeam.CalcShotProbability(match.AwayTeam.GetAbility().Tackling)
		match.AwayTeam.CalcShotProbability(match.HomeTeam.GetAbility().Tackling)
	}
}
