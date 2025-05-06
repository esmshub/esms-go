package flows

import (
	"github.com/esmshub/esms-go/engine/events"
	"github.com/esmshub/esms-go/engine/internal/formulas"
	"github.com/esmshub/esms-go/engine/models"
	"github.com/esmshub/esms-go/engine/pkg/rng"
	"github.com/esmshub/esms-go/engine/pkg/utils"
	"github.com/esmshub/esms-go/engine/types"
	"go.uber.org/zap"
)

var (
	ShotOnTargetProbability            = 5000 // 50%
	ShotOffTargetDeflectionProbability = 975  // 9.75%
	OneOnOneSavedProbability           = 8400 // 84%
)

func ShotFlow(team, oppTeam *models.MatchTeam, flowData map[string]any) []events.Event {
	zap.L().Debug("shot flow", zap.String("team", team.GetName()))

	evts := []events.Event{}

	attacker := utils.MustGetKey[*models.MatchPlayer](flowData, "attacker")
	oppKeeper := oppTeam.GetFirstActiveByPosition(types.PositionGK)
	if oppKeeper == nil {
		zap.L().Panic("no active opposition goalkeeper found", zap.String("team", oppTeam.GetName()))
	}
	flowData["opp_keeper"] = oppKeeper
	// was the shot on target?
	shotOnTarget := rng.Randomp(int(formulas.GetShotOnTargetProbability(attacker.GetMatchAbility().Shooting)))
	goalScored := shotOnTarget && rng.Randomp(int(formulas.GetGoalProbability(attacker.GetMatchAbility().Shooting, oppKeeper.GetMatchAbility().Goalkeeping)))
	flowData["goal_scored"] = goalScored
	oneOnOne := false
	// pre-event
	if assister, ok := flowData["assister"]; ok {
		// chances of a one-on-one situation increases if a goal has been scored
		oneOnOneProb := int(formulas.GetOneOnOneProbability(utils.BoolToInt(goalScored)))
		flowData["one_on_one"] = assister != nil && rng.Randomp(oneOnOneProb)
	}

	if shotOnTarget {
		zap.L().Debug("shot on target", zap.String("attacker", attacker.GetName()))

		evts = append(evts, events.NewEvent(events.ShotOnTargetEventName, flowData))
		evts = append(evts, GoalFlow(team, oppTeam, flowData)...)
	} else {
		zap.L().Debug("shot off target", zap.String("attacker", attacker.GetName()))
		// --- did it hit the post?
		// --- did it hit the bar?
		// --- did it go over the bar?
		// --- did it go past the post?
		if !oneOnOne && rng.Randomp(ShotOffTargetDeflectionProbability) {
			evts = append(evts, events.NewEvent(events.ShotOffTargetDeflectionEventName, flowData))
			evts = append(evts, CornerFlow(team, oppTeam, flowData)...)
		} else {
			evts = append(evts, events.NewEvent(events.ShotOffTargetEventName, flowData))
		}
	}

	return evts
}
