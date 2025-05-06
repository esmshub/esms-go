package flows

import (
	"github.com/esmshub/esms-go/engine/events"
	"github.com/esmshub/esms-go/engine/models"
	"go.uber.org/zap"
)

func FoulFlow(team, oppTeam *models.MatchTeam) []events.Event {
	zap.L().Debug("foul flow", zap.String("team", team.GetName()))
	evts := []events.Event{}
	return evts
}
