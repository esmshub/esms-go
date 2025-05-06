package flows

import (
	"github.com/esmshub/esms-go/engine/events"
	"github.com/esmshub/esms-go/engine/models"
	"go.uber.org/zap"
)

func InjuryFlow(team *models.MatchTeam, opposition *models.MatchTeam) []events.Event {
	zap.L().Debug("injury flow", zap.String("team", team.GetName()))
	evts := []events.Event{}
	return evts
}
