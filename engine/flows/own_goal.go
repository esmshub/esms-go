package flows

import (
	"github.com/esmshub/esms-go/engine/events"
	"github.com/esmshub/esms-go/engine/models"
	"github.com/esmshub/esms-go/engine/pkg/rng"
	"github.com/esmshub/esms-go/engine/types"
	"go.uber.org/zap"
)

var OwnGoalProbability = 40 // 0.4%

func OwnGoalFlow(team, opp *models.MatchTeam) []events.Event {
	zap.L().Debug("own goal flow", zap.String("team", team.GetName()))
	evts := []events.Event{}
	if rng.Randomp(OwnGoalProbability) {
		scorer := getWeightedPlayer(opp.GetActive(), func(p *models.MatchPlayer) int {
			return p.GetBaseAbility().Tackling * 100
		})
		oppKeeper := opp.GetFirstActiveByPosition(types.PositionGK)
		evts = append(evts, events.NewEvent(
			events.OwnGoalScoredEventName,
			map[string]any{
				"scorer":     scorer,
				"opp_keeper": oppKeeper,
			},
		))
	}

	return evts
}
