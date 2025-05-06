package flows

import (
	"github.com/esmshub/esms-go/engine/events"
	"github.com/esmshub/esms-go/engine/internal/formulas"
	"github.com/esmshub/esms-go/engine/models"
	"github.com/esmshub/esms-go/engine/pkg/rng"
	"github.com/esmshub/esms-go/engine/pkg/utils"
	"go.uber.org/zap"
)

var ShotTackledCornerBaseProbability = 975 // 9.75%

func ShotTackleFlow(team, oppTeam *models.MatchTeam, cornerFactor int, flowData map[string]any) []events.Event {
	zap.L().Debug("shot tackle flow", zap.String("team", team.GetName()))
	evts := []events.Event{}
	tackleProb := formulas.CalcShotTackledProbability(
		team.GetAbility().Shooting,
		team.GetAbility().Passing,
		oppTeam.GetAbility().Tackling,
	)
	if rng.Randomp(int(tackleProb)) {
		tackler := utils.MustGetKey[*models.MatchPlayer](flowData, "tackler")
		gotPastDef, hasGotPast := flowData["got_past_defender"]
		if rng.Randomp(ShotTackledCornerBaseProbability * cornerFactor) {
			evts = append(evts, events.NewEvent(events.ShotTackledCornerEventName, flowData))
			evts = append(evts, CornerFlow(team, oppTeam, flowData)...)
		} else if !hasGotPast || gotPastDef.(*models.MatchPlayer).GetName() != tackler.GetName() {
			evts = append(evts, events.NewEvent(events.ShotTackledEventName, flowData))
		} else {
			evts = append(evts, events.NewEvent(events.ShotTackledRecoveryEventName, flowData))
		}
	}

	return evts
}
