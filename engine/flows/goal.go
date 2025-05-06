package flows

import (
	"github.com/esmshub/esms-go/engine/events"
	"github.com/esmshub/esms-go/engine/internal/formulas"
	"github.com/esmshub/esms-go/engine/models"
	"github.com/esmshub/esms-go/engine/pkg/rng"
	"github.com/esmshub/esms-go/engine/pkg/utils"
	"go.uber.org/zap"
)

func GoalFlow(team, oppTeam *models.MatchTeam, flowData map[string]any) []events.Event {
	zap.L().Debug("goal flow", zap.String("team", team.GetName()))
	evts := []events.Event{}

	attacker := utils.MustGetKey[*models.MatchPlayer](flowData, "attacker")
	oppKeeper := utils.MustGetKey[*models.MatchPlayer](flowData, "opp_keeper")
	// this flag is optional, external flows may already ran the goal scored calc
	goalScored, overrideCalc := flowData["goal_scored"]
	// this flag is required but if it's missing we assume it's false
	oneOnOne, _ := flowData["one_on_one"].(bool)

	goalProb := 0
	if !overrideCalc {
		goalProb = int(formulas.GetGoalProbability(
			attacker.GetMatchAbility().Shooting,
			oppKeeper.GetMatchAbility().Goalkeeping,
		))
	}

	if overrideCalc && goalScored.(bool) || !overrideCalc && rng.Randomp(goalProb) {
		// ---- is there a VAR review?
		// ----- was it offside?
		// ----- was it a handball?
		// ----- was it a foul?
		cancelled := rng.Randomp(int(formulas.GoalCancelledProbability))
		flowData["goal_cancelled"] = cancelled
		evts = append(evts, events.NewEvent(events.GoalScoredEventName, flowData))
		if cancelled {
			evts = append(evts, events.NewEvent(events.GoalScoredCancelledEventName, flowData))
		}
	} else if oneOnOne && rng.Randomp(OneOnOneSavedProbability) {
		// ---- was it a corner?
		// ---- was it saved?
		// ----- was it collected?
		// ----- was it a corner?
		// ----- was it parried?
		evts = append(evts, events.NewEvent(events.ShotSavedEventName, flowData))
		if rng.Randomp(int(formulas.CornerFromSaveProbability)) {
			evts = append(evts, CornerFlow(team, oppTeam, flowData)...)
		}
	} else if oneOnOne {
		evts = append(evts, events.NewEvent(events.ShotClearedEventName, flowData))
	} else if rng.Randomp(int(formulas.CornerFromSaveProbability)) {
		// options for shot on target but not a goal
		// --- was it blocked?
		// ---- was it a corner?
		// --- was it saved?
		// ---- was it collected?
		// ---- was it a corner?
		// ---- was it parried?
		evts = append(evts, events.NewEvent(events.ShotSavedCornerEventName, flowData))
		evts = append(evts, CornerFlow(team, oppTeam, flowData)...)
	} else {
		evts = append(evts, events.NewEvent(events.ShotSavedEventName, flowData))
	}

	return evts
}
