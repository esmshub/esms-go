package models

type Player struct {
	Name     string
	Age      int
	Nat      string
	Position string
	Ability  *PlayerAbilities
	Stats    *PlayerStats
}

func (p *Player) GetIsInjured() bool {
	return p.Stats.WeeksInjured > 0
}

func (p *Player) GetIsSuspended() bool {
	return p.Stats.GamesSuspended > 0
}

type PlayerAbilityPoints struct {
	GoalkeepingAbs int
	TacklingAbs    int
	PassingAbs     int
	ShootingAbs    int
}

type PlayerAbilities struct {
	Goalkeeping int
	Tackling    int
	Passing     int
	Shooting    int
	Aggression  int
	PlayerAbilityPoints
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

type PlayerGameStats struct {
	MinutesPlayed int
	IsMom         bool
	Saves         int // GK only
	KeyTackles    int
	KeyPasses     int
	Assists       int
	Shots         int
	Goals         int
	Fouls         int
	IsCautioned   bool
	IsSentOff     bool
	IsInjured     bool
	IsSuspended   bool
	IsSubbed      bool
}

type PlayerPosition struct {
	Position    string
	Name        string
	BaseAbility *PlayerAbilities
	Ability     *PlayerAbilities
	Stats       *PlayerGameStats
	IsActive    bool
	Fatigue     float64
}
