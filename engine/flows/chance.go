package flows

import (
	"slices"
	"strings"

	"github.com/esmshub/esms-go/engine/events"
	"github.com/esmshub/esms-go/engine/models"
	"github.com/esmshub/esms-go/engine/pkg/rng"
	"go.uber.org/zap"
)

var (
	ChanceAssistedProbability              = 7000 // 70%
	ChanceAssistedBeatsDefenderProbability = 700  // 7%
	ChanceBeatsDefenderProbability         = 4000 // 40%
)

func matchEventName(name string) func(events.Event) bool {
	return func(e events.Event) bool {
		return strings.EqualFold(name, e.GetName())
	}
}

func ChanceFlow(team, oppTeam *models.MatchTeam) []events.Event {
	zap.L().Debug("chance flow", zap.String("team", team.GetName()))
	evts := []events.Event{}
	data := map[string]any{}
	if rng.Randomp(int(team.GetShotProbability())) {
		zap.L().Debug("shooting opportunity")
		// was there an own goal?
		evts = append(evts, OwnGoalFlow(team, oppTeam)...)
		if !slices.ContainsFunc(evts, matchEventName(events.OwnGoalScoredEventName)) {
			// choose attacker
			attacker := getAttacker(team)
			data["attacker"] = attacker

			// was there an assister?
			var assister *models.MatchPlayer
			if rng.Randomp(ChanceAssistedProbability) {
				assister = getAssister(team)
				if assister.GetName() == attacker.GetName() {
					assister = nil
				} else {
					zap.L().Debug("chance assisted", zap.String("assister", assister.GetName()))
					data["assister"] = assister
				}
			}

			if assister != nil {
				if rng.Randomp(ChanceAssistedBeatsDefenderProbability) {
					data["got_past_defender"] = getDefender(oppTeam)
					evts = append(evts, events.NewEvent(events.AssistedChanceBeatsDefenderEventName, data))
				} else {
					evts = append(evts, events.NewEvent(events.AssistedChanceEventName, data))
				}
			} else {
				if rng.Randomp(ChanceBeatsDefenderProbability) {
					data["got_past_defender"] = getDefender(oppTeam)
					evts = append(evts, events.NewEvent(events.ChanceBeatsDefenderEventName, data))
				} else {
					evts = append(evts, events.NewEvent(events.ChanceEventName, data))
				}
			}

			data["tackler"] = getDefender(oppTeam)
			// was the chance tackled?
			eventCountBefore := len(evts)
			evts = append(evts, ShotTackleFlow(team, oppTeam, 1, data)...)
			if eventCountBefore == len(evts) {
				// chance was not tackled, there was an attempt at goal
				evts = append(evts, ShotFlow(team, oppTeam, data)...)
			}
		}
	}

	return evts
}
