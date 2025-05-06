package models

import "sync"

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
	MinutesPlayed  int
	IsMom          bool
	Saves          int // GK only
	KeyTackles     int
	KeyPasses      int
	Assists        int
	ShotsOnTarget  int
	ShotsOffTarget int
	Goals          []int
	Fouls          int
	Conceded       int
	OwnGoals       []int
	IsCautioned    bool
	IsSentOff      bool
	IsInjured      bool
	// IsSuspended    bool
	IsSubbed bool
}

type MatchPlayer struct {
	mut       *sync.RWMutex
	condition float64
	player    *Player
	stats     *PlayerGameStats
	ability   *PlayerAbilities
	IsActive  bool
	IsSub     bool
}

func (p *MatchPlayer) GetName() string {
	return p.player.Name
}

func (p *MatchPlayer) GetPosition() string {
	return p.player.Position
}

func (p *MatchPlayer) GetMatchAbility() PlayerAbilities {
	p.mut.RLock()
	defer p.mut.RUnlock()

	return *p.ability
}

func (p *MatchPlayer) GetBaseAbility() PlayerAbilities {
	p.mut.RLock()
	defer p.mut.RUnlock()

	return *p.player.Ability
}

func (p *MatchPlayer) GetCondition() float64 {
	p.mut.RLock()
	defer p.mut.RUnlock()

	return p.condition
}

func (p *MatchPlayer) SetCondition(condition float64) {
	p.mut.Lock()
	defer p.mut.Unlock()

	p.condition = condition
}

func (p *MatchPlayer) SetMatchAbility(ability *PlayerAbilities) {
	p.mut.Lock()
	defer p.mut.Unlock()

	p.ability = ability
}

func (p *MatchPlayer) AddMinute() {
	p.mut.Lock()
	defer p.mut.Unlock()

	p.stats.MinutesPlayed++
}

func (p *MatchPlayer) AddGoal(minute int) {
	p.mut.Lock()
	defer p.mut.Unlock()

	p.stats.Goals = append(p.stats.Goals, minute)
}

func (p *MatchPlayer) AddAssist() {
	p.mut.Lock()
	defer p.mut.Unlock()

	p.stats.Assists++
}

func (p *MatchPlayer) AddSave() {
	p.mut.Lock()
	defer p.mut.Unlock()

	p.stats.Saves++
}

func (p *MatchPlayer) AddKeyTackle() {
	p.mut.Lock()
	defer p.mut.Unlock()

	p.stats.KeyTackles++
}

func (p *MatchPlayer) AddKeyPass() {
	p.mut.Lock()
	defer p.mut.Unlock()

	p.stats.KeyPasses++
}

func (p *MatchPlayer) AddShotOnTarget() {
	p.mut.Lock()
	defer p.mut.Unlock()

	p.stats.ShotsOnTarget++
}

func (p *MatchPlayer) AddShotOffTarget() {
	p.mut.Lock()
	defer p.mut.Unlock()

	p.stats.ShotsOffTarget++
}

func (p *MatchPlayer) AddFoul() {
	p.mut.Lock()
	defer p.mut.Unlock()

	p.stats.Fouls++
}

func (p *MatchPlayer) AddConceded() {
	p.mut.Lock()
	defer p.mut.Unlock()

	p.stats.Conceded++
}

func (p *MatchPlayer) AddOwnGoal(minute int) {
	p.mut.Lock()
	defer p.mut.Unlock()

	p.stats.OwnGoals = append(p.stats.OwnGoals, minute)
}

func (p *MatchPlayer) SetIsCautioned(isCautioned bool) {
	p.mut.Lock()
	defer p.mut.Unlock()

	p.stats.IsCautioned = isCautioned
}

func (p *MatchPlayer) SetIsSentOff(isSentOff bool) {
	p.mut.Lock()
	defer p.mut.Unlock()

	p.stats.IsSentOff = isSentOff
}

func (p *MatchPlayer) SetIsInjured(isInjured bool) {
	p.mut.Lock()
	defer p.mut.Unlock()

	p.stats.IsInjured = isInjured
}

func (p *MatchPlayer) SetIsSubbed(isSubbed bool) {
	p.mut.Lock()
	defer p.mut.Unlock()

	p.stats.IsSubbed = isSubbed
}

func (p *MatchPlayer) GetIsInjured() bool {
	p.mut.RLock()
	defer p.mut.RUnlock()

	return p.player.GetIsInjured() || p.stats.IsInjured
}

func (p *MatchPlayer) GetIsSuspended() bool {
	p.mut.RLock()
	defer p.mut.RUnlock()

	return p.player.GetIsSuspended()
}

func (p *MatchPlayer) GetStats() PlayerGameStats {
	p.mut.RLock()
	defer p.mut.RUnlock()

	return *p.stats
}

func NewMatchPlayer(player *Player) *MatchPlayer {
	return &MatchPlayer{
		mut:       &sync.RWMutex{},
		player:    player,
		condition: 1.0,
		// copy base abilities
		ability: &PlayerAbilities{
			Goalkeeping: player.Ability.Goalkeeping,
			Tackling:    player.Ability.Tackling,
			Passing:     player.Ability.Passing,
			Shooting:    player.Ability.Shooting,
			Aggression:  player.Ability.Aggression,
		},
		stats: &PlayerGameStats{},
	}
}
