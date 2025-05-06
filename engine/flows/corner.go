package flows

import (
	"github.com/esmshub/esms-go/engine/events"
	"github.com/esmshub/esms-go/engine/models"
	"github.com/esmshub/esms-go/engine/pkg/rng"
	"github.com/esmshub/esms-go/engine/types"
	"go.uber.org/zap"
)

// var OwnGoalProbability = 40 // 0.4%
var CornerKeeperCatchProbability = 2250 // 22.5%
var CornerClearedProbability = 3740     // 37.4%

func CornerFlow(team, oppTeam *models.MatchTeam, flowData map[string]any) []events.Event {
	zap.L().Debug("corner flow", zap.String("team", team.GetName()))
	evts := []events.Event{}

	cornerTaker := getWeightedPlayer(team.GetActive(), func(p *models.MatchPlayer) int {
		return p.GetMatchAbility().Passing * 100
	})
	flowData["corner_taker"] = cornerTaker
	oppKeeper := oppTeam.GetFirstActiveByPosition(types.PositionGK)
	keeperCatchesProb := CornerKeeperCatchProbability + (oppKeeper.GetBaseAbility().Goalkeeping * 900) - (cornerTaker.GetBaseAbility().Passing * 1150)
	defenderClearsProb := CornerClearedProbability * (oppTeam.GetAbility().Tackling * 3) / ((team.GetAbility().Passing * 2) + team.GetAbility().Shooting)

	if rng.Randomp(keeperCatchesProb) {
		flowData["opp_keeper"] = oppKeeper
		evts = append(evts, events.NewEvent(events.CornerCaughtEventName, flowData))
	} else if rng.Randomp(defenderClearsProb) {
		flowData["corner_defender"] = getWeightedPlayer(oppTeam.GetActive(), func(p *models.MatchPlayer) int {
			return p.GetMatchAbility().Passing * 100
		})
		evts = append(evts, events.NewEvent(events.CornerClearedEventName, flowData))
	} else {
		// shot comes from corner
		flowData["attacker"] = getWeightedPlayer(team.GetActive(), func(p *models.MatchPlayer) int {
			return (p.GetBaseAbility().Shooting + 10) * 100
		})
		evts = append(evts, events.NewEvent(events.CornerShotEventName, flowData))
		evts = append(evts, ShotFlow(team, oppTeam, flowData)...)
	}

	return evts
}
