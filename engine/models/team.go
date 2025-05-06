package models

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/esmshub/esms-go/engine/internal/formulas"
	"github.com/esmshub/esms-go/engine/pkg/utils"
	"github.com/esmshub/esms-go/engine/types"
)

type TeamStats struct {
	Goals          []MatchGoal
	Penalties      []MatchGoal
	Possession     int
	Tackles        int
	Passes         int
	ShotsOnTarget  int
	ShotsOffTarget int
	YellowCards    int
	RedCards       int
	// FoulsCommitted int
	// Substitutions  int
	// Possession     float64
}

type MatchGoal struct {
	Player *MatchPlayer
	Minute int
}

type TeamConfig struct {
	Name            string
	Code            string
	ManagerName     string
	StadiumName     string
	StadiumCapacity int
	Tactic          string
	Players         []*MatchPlayer
	Roster          []*Player
	PlayerRoles     map[string]*MatchPlayer
	Conditionals    []*Conditional
	TeamAbility     *PlayerAbilities
	Injuries        int
}

func (t *TeamConfig) GetStarters() []*MatchPlayer {
	return utils.FilterFunc(t.Players, func(p *MatchPlayer) bool {
		return !p.IsSub
	})
}

func (t *TeamConfig) GetSubs() []*MatchPlayer {
	return utils.FilterFunc(t.Players, func(p *MatchPlayer) bool {
		return p.IsSub
	})
}

func (t *TeamConfig) GetInjuredPlayers() []*Player {
	var injured []*Player
	for _, p := range t.Roster {
		if p.GetIsInjured() {
			injured = append(injured, p)
		}
	}
	return injured
}

func (t *TeamConfig) GetSuspendedPlayers() []*Player {
	var suspended []*Player
	for _, p := range t.Roster {
		if p.GetIsSuspended() {
			suspended = append(suspended, p)
		}
	}
	return suspended
}

func (t *TeamConfig) GetFormation() string {
	positions := make(map[string]int)
	for _, p := range t.GetStarters() {
		positions[p.GetPosition()]++
	}
	defMids := positions[types.PositionDM]
	mids := positions[types.PositionMF]
	attMids := positions[types.PositionAM]

	formationStr := fmt.Sprintf("%d-%d-%d-%d-%d %s", positions[types.PositionDF], defMids, mids, attMids, positions[types.PositionFW], t.Tactic)
	return strings.ReplaceAll(formationStr, "-0", "")
}

func (t *TeamConfig) Accept(visitor Visitor) {
	visitor.VisitTeam(t)
}

type MatchTeam struct {
	mut             *sync.RWMutex
	config          *TeamConfig
	stats           *TeamStats
	ability         *PlayerAbilities
	shotProbability float64
}

func (t *MatchTeam) GetShortName() string {
	t.mut.RLock()
	defer t.mut.RUnlock()

	return t.config.Code
}

func (t *MatchTeam) GetName() string {
	t.mut.RLock()
	defer t.mut.RUnlock()

	return t.config.Name
}

func (t *MatchTeam) GetManagerName() string {
	t.mut.RLock()
	defer t.mut.RUnlock()

	return t.config.ManagerName
}

func (t *MatchTeam) GetStadiumName() string {
	t.mut.RLock()
	defer t.mut.RUnlock()

	return t.config.StadiumName
}

func (t *MatchTeam) GetFormation() string {
	return t.config.GetFormation()
}

func (t *MatchTeam) GetTactic() string {
	return t.config.Tactic
}

func (t *MatchTeam) GetStarters() []*MatchPlayer {
	return t.config.GetStarters()
}

func (t *MatchTeam) GetSubs() []*MatchPlayer {
	return t.config.GetSubs()
}

func (t *MatchTeam) GetLineup() []*MatchPlayer {
	t.mut.RLock()
	defer t.mut.RUnlock()

	return t.config.Players
}

func (t *MatchTeam) GetAbility() PlayerAbilities {
	t.mut.RLock()
	defer t.mut.RUnlock()

	return *t.ability
}

func (t *MatchTeam) GetStats() TeamStats {
	t.mut.RLock()
	defer t.mut.RUnlock()

	return *t.stats
}

func (t *MatchTeam) GetConditionals() []*Conditional {
	t.mut.RLock()
	defer t.mut.RUnlock()

	return t.config.Conditionals
}

func (t *MatchTeam) SetAbility(ability *PlayerAbilities) {
	t.mut.Lock()
	defer t.mut.Unlock()

	t.ability = ability
}

func (t *MatchTeam) AddGoal(player *MatchPlayer, minute int) {
	t.mut.Lock()
	defer t.mut.Unlock()

	t.stats.Goals = append(t.stats.Goals, MatchGoal{
		Player: player,
		Minute: minute,
	})
}

func (t *MatchTeam) AddTackle() {
	t.mut.Lock()
	defer t.mut.Unlock()

	t.stats.Tackles++
}

func (t *MatchTeam) AddPass() {
	t.mut.Lock()
	defer t.mut.Unlock()

	t.stats.Passes++
}

func (t *MatchTeam) AddShotOnTarget() {
	t.mut.Lock()
	defer t.mut.Unlock()

	t.stats.ShotsOnTarget++
}

func (t *MatchTeam) AddShotOffTarget() {
	t.mut.Lock()
	defer t.mut.Unlock()

	t.stats.ShotsOffTarget++
}

func (t *MatchTeam) SetPossession(possession int) {
	t.mut.Lock()
	defer t.mut.Unlock()

	t.stats.Possession = int(math.Min(math.Max(float64(possession), 0), 100))
}

func (t *MatchTeam) GetShotProbability() float64 {
	t.mut.RLock()
	defer t.mut.RUnlock()

	return t.shotProbability
}

func (t *MatchTeam) CalcShotProbability(oppositionTackling int) {
	t.mut.Lock()
	defer t.mut.Unlock()

	t.shotProbability = formulas.CalcShotProbability(
		t.ability.Aggression,
		t.ability.Shooting,
		t.ability.Passing,
		oppositionTackling,
	)
}

func (t *MatchTeam) GetActive() []*MatchPlayer {
	t.mut.RLock()
	defer t.mut.RUnlock()

	return utils.FilterFunc(t.config.Players, func(p *MatchPlayer) bool {
		return p.IsActive
	})
}

func (t *MatchTeam) GetActiveByPosition(position string) []*MatchPlayer {
	t.mut.RLock()
	defer t.mut.RUnlock()

	return utils.FilterFunc(t.config.Players, func(p *MatchPlayer) bool {
		return p.IsActive && p.GetPosition() == position
	})
}

func (t *MatchTeam) GetFirstActiveByPosition(position string) *MatchPlayer {
	t.mut.RLock()
	defer t.mut.RUnlock()

	return utils.FindFunc(t.config.Players, func(p *MatchPlayer) bool {
		return p.IsActive && p.GetPosition() == position
	})
}

func NewMatchTeam(config *TeamConfig) *MatchTeam {
	return &MatchTeam{
		mut:    &sync.RWMutex{},
		config: config,
		stats: &TeamStats{
			Goals:      []MatchGoal{},
			Possession: 0,
			Tackles:    0,
		},
		ability: &PlayerAbilities{},
	}
}
