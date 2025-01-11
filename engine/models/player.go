package models

type Player struct {
	Name      string
	Age       int
	Nat       string
	Position  string
	Abilities *PlayerAbilities
	Stats     *PlayerStats
}

func (p *Player) GetIsInjured() bool {
	return p.Stats.WeeksInjured > 0
}

func (p *Player) GetIsSuspended() bool {
	return p.Stats.GamesSuspended > 0
}

type PlayerAbilities struct {
	Goalkeeping    int
	Tackling       int
	Passing        int
	Shooting       int
	Aggression     int
	GoalkeepingAbs int
	TacklingAbs    int
	PassingAbs     int
	ShootingAbs    int
}

type PlayerStats struct {
	GamesStarted       int
	GamesSubbed        int
	MinutesPlayed      int
	MomAwards          int
	Saves              int // GK only
	GoalsConceded      int // GK only
	KeyTackles         int
	KeyPasses          int
	Shots              int
	Goals              int
	Assists            int
	DisciplinaryPoints int
	WeeksInjured       int
	GamesSuspended     int
}

type PlayerPosition struct {
	Position string
	Name     string
	Stats    *PlayerStats
}
