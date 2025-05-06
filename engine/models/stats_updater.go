package models

import (
	"fmt"

	"github.com/esmshub/esms-go/engine/events"
	"github.com/esmshub/esms-go/engine/pkg/utils"
	"go.uber.org/zap"
)

type MatchStatsUpdater struct{}

func (m *MatchStatsUpdater) Update(event string, subject Subject) {
	match, ok := subject.(*Match)
	if !ok {
		zap.L().DPanic("invalid subject type", zap.String("type", fmt.Sprintf("%T", subject)))
	}

	if event == events.MinuteElapsedEventName {
		utils.EachFunc(match.GetTeams(), m.Visit)
	}
}

func (m *MatchStatsUpdater) Visit(team *MatchTeam) {
	for _, p := range team.GetActive() {
		p.AddMinute()
		if p.GetName() == "FW2" {
			zap.L().Info("player minute", zap.String("name", p.GetName()), zap.Int("minute", p.GetStats().MinutesPlayed))
		}
	}
}
