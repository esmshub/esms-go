package types

const (
	MatchTypeLeague = iota
	MatchTypeCup
	MatchTypeFriendly
)

const (
	AggressionAction       string = "AGG"
	ChangeAggressionAction string = "CHANGEAGG"
	ChangePositionAction   string = "CHANGEPOS"
	SubstituteAction       string = "SUB"
	ChangeTacticAction     string = "TACTIC"
	InjuryEvent            string = "INJURED"
	MinuteEvent            string = "MIN"
	RedCardEvent           string = "RED"
	YellowCardEvent        string = "YELLOW"
	ShotEvent              string = "SHOTS"
	ScoreEvent             string = "SCORE"
	TacticCounter                 = "C"
	TacticAttacking               = "A"
	TacticEuropean                = "E"
	TacticDefensive               = "D"
	TacticNormal                  = "N"
	TacticLong                    = "L"
	TacticPassing                 = "P"
	TacticT                       = "T"
	TacticCounterName             = "Counter Attack"
	TacticAttackingName           = "Attacking"
	TacticDefensiveName           = "Defensive"
	TacticNormalName              = "Normal"
	TacticLongName                = "Long Ball"
	TacticEuropeanName            = "European"
	TacticPassingName             = "Passing"
	TacticTName                   = "T Tactic"
	PositionGK                    = "GK"
	PositionDF                    = "DF"
	PositionDM                    = "DM"
	PositionMF                    = "MF"
	PositionAM                    = "AM"
	PositionFW                    = "FW"
	RolePenaltyTaker              = "PK"
)

var (
	ValidPositions = []string{
		PositionGK,
		PositionDF,
		PositionDM,
		PositionMF,
		PositionAM,
		PositionFW,
	}
	ValidTactics = []string{
		TacticCounter,
		TacticAttacking,
		TacticEuropean,
		TacticDefensive,
		TacticNormal,
		TacticLong,
		TacticPassing,
		TacticT,
	}
	TacticNames = map[string]string{
		TacticCounter:   TacticCounterName,
		TacticAttacking: TacticAttackingName,
		TacticEuropean:  TacticEuropeanName,
		TacticDefensive: TacticDefensiveName,
		TacticNormal:    TacticNormalName,
		TacticLong:      TacticLongName,
		TacticPassing:   TacticPassingName,
		TacticT:         TacticTName,
	}
	ValidConditionalEvents = []string{
		InjuryEvent,
		MinuteEvent,
		RedCardEvent,
		YellowCardEvent,
		ShotEvent,
		ScoreEvent,
	}
)

type MATCHTYPE int

type ReadableOptions interface {
	Get(string) any
}
