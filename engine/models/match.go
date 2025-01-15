package models

import "github.com/esmshub/esms-go/engine/common"

type Referee struct {
	Name string
	Nat  string
}

type Match struct {
	HomeTeam   *TeamConfig
	AwayTeam   *TeamConfig
	Referee    *Referee
	Commentary common.CommentaryProvider
}

type MatchResult struct {
	HomeTeam *TeamConfig
	AwayTeam *TeamConfig
	Referee  *Referee
	RngSeed  uint64
}
