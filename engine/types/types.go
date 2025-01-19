package types

const (
	MATCHTYPE_LEAGUE = iota
	MATCHTYPE_CUP
	MATCHTYPE_FRIENDLY
)

const (
	AGG_ACTION            string = "AGG"
	CHANGEAGG_ACTION      string = "CHANGEAGG"
	CHANGEPOS_ACTION      string = "CHANGEPOS"
	SUB_ACTION            string = "SUB"
	TACTIC_ACTION         string = "TACTIC"
	INJ_EVENT             string = "INJURED"
	MIN_EVENT             string = "MIN"
	REDCARD_EVENT         string = "RED"
	YELLOWCARD_EVENT      string = "YELLOW"
	SHOT_EVENT            string = "SHOTS"
	SCORE_EVENT           string = "SCORE"
	TACTIC_COUNTER               = "C"
	TACTIC_ATTACKING             = "A"
	TACTIC_EUROPEAN              = "E"
	TACTIC_DEFENSIVE             = "D"
	TACTIC_NORMAL                = "N"
	TACTIC_LONG                  = "L"
	TACTIC_PASSING               = "P"
	TACTIC_COUNTER_NAME          = "Counter Attack"
	TACTIC_ATTACKING_NAME        = "Attacking"
	TACTIC_DEFENSIVE_NAME        = "Defensive"
	TACTIC_NORMAL_NAME           = "Normal"
	TACTIC_LONG_NAME             = "Long Ball"
	TACTIC_EUROPEAN_NAME         = "European"
	TACTIC_PASSING_NAME          = "Passing"
	POSITION_GK                  = "GK"
	POSITION_DF                  = "DF"
	POSITION_DM                  = "DM"
	POSITION_MF                  = "MF"
	POSITION_AM                  = "AM"
	POSITION_FW                  = "FW"
)

var (
	ValidPositions = []string{
		POSITION_GK,
		POSITION_DF,
		POSITION_DM,
		POSITION_MF,
		POSITION_AM,
		POSITION_FW,
	}
	ValidTactics = []string{
		TACTIC_COUNTER,
		TACTIC_ATTACKING,
		TACTIC_EUROPEAN,
		TACTIC_DEFENSIVE,
		TACTIC_NORMAL,
		TACTIC_LONG,
		TACTIC_PASSING,
	}
	TacticNames = map[string]string{
		TACTIC_COUNTER:   TACTIC_COUNTER_NAME,
		TACTIC_ATTACKING: TACTIC_ATTACKING_NAME,
		TACTIC_EUROPEAN:  TACTIC_EUROPEAN_NAME,
		TACTIC_DEFENSIVE: TACTIC_DEFENSIVE_NAME,
		TACTIC_NORMAL:    TACTIC_NORMAL_NAME,
		TACTIC_LONG:      TACTIC_LONG_NAME,
		TACTIC_PASSING:   TACTIC_PASSING_NAME,
	}
	ValidConditionalEvents = []string{
		INJ_EVENT,
		MIN_EVENT,
		REDCARD_EVENT,
		YELLOWCARD_EVENT,
		SHOT_EVENT,
		SCORE_EVENT,
	}
)

type MATCHTYPE int

type ReadableOptions interface {
	Get(string) any
}
