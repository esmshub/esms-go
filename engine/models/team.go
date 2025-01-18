package models

import (
	"github.com/esmshub/esms-go/engine/internal/formulas"
)

type TeamConfig struct {
	Name            string
	Code            string
	ManagerName     string
	StadiumName     string
	StadiumCapacity int
	Formation       string
	Tactic          string
	Lineup          []*PlayerPosition
	Subs            []*PlayerPosition
	Roster          []*Player
	PlayerRoles     map[string]*PlayerPosition
	Conditionals    []*Conditional
	TeamAbility     *PlayerAbilities
	ShotProbability float64
}

func (r *TeamConfig) GetShotProbability(oppositionAbility *PlayerAbilities) float64 {
	return formulas.CalcShotProbability(
		float64(r.TeamAbility.Aggression),
		float64(r.TeamAbility.Shooting),
		float64(r.TeamAbility.Passing),
		float64(oppositionAbility.Tackling),
	)
}

func (r *TeamConfig) GetInjuredPlayers() []*Player {
	var injured []*Player
	for _, p := range r.Roster {
		if p.GetIsInjured() {
			injured = append(injured, p)
		}
	}
	return injured
}

func (r *TeamConfig) GetSuspendedPlayers() []*Player {
	var suspended []*Player
	for _, p := range r.Roster {
		if p.GetIsSuspended() {
			suspended = append(suspended, p)
		}
	}
	return suspended
}

func (r *TeamConfig) DecreaseFatigue(minutes int, teams []*TeamConfig) {
	r.IncreasePlayerFatigue(-minutes, teams)
}

func (r *TeamConfig) IncreasePlayerFatigue(minutes int, teams []*TeamConfig) {
	for _, t := range teams {
		players := append(t.Lineup, t.Subs...)
		for _, attrs := range players {
			if attrs.IsActive {
				// fmt.Printf("%s fatigue adjusted from %f", attrs.BaseStats.Name, attrs.Fatigue)
				attrs.Fatigue = formulas.CalcFatigue(attrs.Fatigue, minutes)
				// fmt.Printf(" to %f\n", attrs.Fatigue)
			}
		}
	}
}
