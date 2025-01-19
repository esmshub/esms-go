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
	HomeTeam   *TeamConfig
	AwayTeam   *TeamConfig
	Possession [2]int
	Referee    *Referee
	RngSeed    uint64
}

type TeamStats struct {
	Goals      int
	Possession int
	// ShotsOnTarget  int
	// ShotsOffTarget int
	// FoulsCommitted int
	// Substitutions  int
	// Possession     float64
}

type MatchStats struct {
	HomeStats TeamStats
	AwayStats TeamStats
}
