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
