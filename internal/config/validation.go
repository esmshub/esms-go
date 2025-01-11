package config

import (
	"slices"
	"strings"
)

var (
	POSITION_GK      = "GK"
	POSITION_DF      = "DF"
	POSITION_DM      = "DM"
	POSITION_MF      = "MF"
	POSITION_AM      = "AM"
	POSITION_FW      = "FW"
	TACTIC_COUNTER   = "Counter Attack"
	TACTIC_ATTACKING = "Attacking"
	TACTIC_DEFENSIVE = "Defensive"
	TACTIC_NORMAL    = "Normal"
	TACTIC_LONG      = "Long Ball"
	TACTIC_EUROPEAN  = "European"
	TACTIC_PASSING   = "Passing"
)

var TacticNames = map[string]string{
	"C": TACTIC_COUNTER,
	"A": TACTIC_ATTACKING,
	"E": TACTIC_EUROPEAN,
	"D": TACTIC_DEFENSIVE,
	"N": TACTIC_NORMAL,
	"L": TACTIC_LONG,
	"P": TACTIC_PASSING,
}

var ValidPositions = []string{
	POSITION_GK,
	POSITION_DF,
	POSITION_DM,
	POSITION_MF,
	POSITION_AM,
	POSITION_FW,
}

func IsValidPosition(position string) bool {
	return slices.ContainsFunc(ValidPositions, func(pos string) bool {
		return strings.EqualFold(position, pos)
	})
}
