package models

import "github.com/esmshub/esms-go/engine/pkg/utils"

type Validator func(TeamConfig) error

type TeamConfig struct {
	validators   []Validator
	Name         string
	Formation    string
	Tactic       string
	Lineup       []*Player
	Subs         []*Player
	Roster       []*Player
	PlayerRoles  map[string]*Player
	Conditionals []*Conditional
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
