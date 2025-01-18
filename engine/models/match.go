package models

type Referee struct {
	Name string
	Nat  string
}

type Match struct {
	HomeTeam   *TeamConfig
	AwayTeam   *TeamConfig
	Referee    *Referee
	Commentary CommentaryProvider
}

type MatchResult struct {
	HomeTeam *TeamConfig
	AwayTeam *TeamConfig
	Referee  *Referee
	RngSeed  uint64
}
