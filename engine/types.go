package engine

import (
	"github.com/esmshub/esms-go/engine/commentary"
	"github.com/esmshub/esms-go/engine/models"
	"github.com/esmshub/esms-go/engine/types"
)

type Options struct {
	MatchType          types.MATCHTYPE
	TacticsMatrix      *models.TacticsMatrix
	CommentaryProvider commentary.CommentaryProvider
	RngSeed            uint64
	AppConfig          map[string]any
	Bonuses            map[string]any
}

const MaxSkillValue = 20
const MaxTeamSkillValue = 220

const PossessionTacklingWeight = 0.3
const PossessionPassingWeight = 0.55
const PossessionShootingWeight = 0.15

const MinsPerHalf = 45
const DefAggressionLevel = 10

// var EVENT_FLOWS = map[string]flows.Flow{
// 	types.SHOT_EVENT:     flows.ShotFlow,
// 	types.TA:     flows.FoulFlow,
// 	types.TACKLE_EVENT:   flows.TackleFlow,
// }
