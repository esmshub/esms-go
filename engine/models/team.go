package models

import "github.com/esmshub/esms-go/engine/pkg/utils"

type Validator func(TeamConfig) error

type TeamConfig struct {
	validators      []Validator
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
}

func (r *TeamConfig) Validate() []error {
	return utils.Reduce(r.validators, func(results []error, v Validator) []error {
		return append(results, v(*r))
	}, []error{})
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
