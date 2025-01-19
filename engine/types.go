package engine

import (
	"github.com/esmshub/esms-go/engine/models"
	"github.com/esmshub/esms-go/engine/types"
)

type Options struct {
	MatchType     types.MATCHTYPE
	TacticsMatrix *models.TacticsMatrix
	RngSeed       uint64
	AppConfig     map[string]any
}

const MAX_SKILL_VALUE = 20
const MAX_TEAMSKILL_VALUE = 220

const POS_TK_WEIGHT = 0.3
const POS_PS_WEIGHT = 0.55
const POS_SH_WEIGHT = 0.15

const MINS_PER_HALF = 45
const DEFAULT_AGG_LEVEL = 10
