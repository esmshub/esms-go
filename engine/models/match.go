package models

import (
	"sync"

	"github.com/esmshub/esms-go/engine/events"
	"github.com/esmshub/esms-go/engine/internal/defaults"
	"github.com/esmshub/esms-go/engine/pkg/rng"
	"github.com/esmshub/esms-go/engine/pkg/utils"
	"go.uber.org/zap"
)

type Referee struct {
	Name string
	Nat  string
}

type Match struct {
	mut        *sync.RWMutex
	minute     int
	injuryTime int
	HomeTeam   *MatchTeam
	AwayTeam   *MatchTeam
	Referee    *Referee
	*SubjectImpl
}

type MatchResult struct {
	HomeTeam *MatchTeam
	AwayTeam *MatchTeam
	Referee  *Referee
	RngSeed  uint64
}

func (m *MatchResult) IsWinner(team *MatchTeam) bool {
	if m.HomeTeam.GetShortName() == team.GetShortName() {
		return len(m.HomeTeam.GetStats().Goals) > len(m.AwayTeam.GetStats().Goals)
	} else if m.AwayTeam.GetShortName() == team.GetShortName() {
		return len(m.AwayTeam.GetStats().Goals) > len(m.HomeTeam.GetStats().Goals)
	} else {
		zap.L().DPanic("team not found", zap.String("team", team.GetShortName()))
	}

	return false
}

func (m *MatchResult) HasCleanSheet(team *MatchTeam) bool {
	if m.HomeTeam.GetShortName() == team.GetShortName() {
		return len(m.AwayTeam.GetStats().Goals) == 0
	} else if m.AwayTeam.GetShortName() == team.GetShortName() {
		return len(m.HomeTeam.GetStats().Goals) == 0
	}

	return false
}

type MatchStats struct {
	HomeStats  TeamStats
	AwayStats  TeamStats
	InjuryTime int
}

func (mt *Match) CalculateInjuryTime() {
	mt.mut.Lock()
	mt.injuryTime = rng.RandomRange(0, 5)
	mt.mut.Unlock()

	mt.notify(events.InjuryTimeAddedEventName)
}

func (mt *Match) GetInjuryTime() int {
	mt.mut.RLock()
	defer mt.mut.RUnlock()

	return mt.injuryTime
}

func (mt *Match) IncrementMinute() {
	mt.mut.Lock()
	mt.minute++
	mt.mut.Unlock()

	mt.notify(events.MinuteElapsedEventName)
}

// Calculates how much injury time to add.
//
// Takes into account substitutions, injuries and fouls (by both teams)
func (m *Match) calcMaxInjuryTime() {
	/*
		for consideration:
		- The assessment of apparently injured players
		- The removal from the field of injured players
		- Substitutions
		- Perceived time wasting by players
		- Red or yellow cards being issued
		- Delays for VAR checks
		- Drinks breaks in hotter venues
	*/
	// teams := []*models.TeamConfig{match.HomeTeam, match.AwayTeam}
	// subCount := utils.Reduce(teams, func(acc int, t *models.TeamConfig) int {
	// 	return acc + utils.CountFunc(t.Lineup, func(p *models.PlayerPosition) bool {
	// 		return p.Stats.IsSubbed
	// 	}) + utils.CountFunc(t.Subs, func(p *models.PlayerPosition) bool {
	// 		return p.Stats.IsSubbed
	// 	})
	// }, 0)
	// injuryCount := utils.Reduce(teams, func(acc int, t *models.TeamConfig) int {
	// 	return acc + utils.CountFunc(t.Lineup, func(p *models.PlayerPosition) bool {
	// 		return p.Stats.IsInjured
	// 	}) + utils.CountFunc(t.Subs, func(p *models.PlayerPosition) bool {
	// 		return p.Stats.IsInjured
	// 	})
	// }, 0)
	// foulCount := utils.Reduce(teams, func(acc int, t *models.TeamConfig) int {
	// 	return acc + utils.SumFunc(t.Lineup, func(p *models.PlayerPosition) int {
	// 		return p.Stats.Fouls
	// 	}) + utils.SumFunc(t.Subs, func(p *models.PlayerPosition) int {
	// 		return p.Stats.Fouls
	// 	})
	// }, 0)

	// return int(math.Ceil(float64(subCount+injuryCount+foulCount) * 0.5))

	// TODO: revisit true calc above
	// if len(m.injuryTime) < 2 {
	// 	result := rng.RandomRange(1, 5)
	// 	m.injuryTime = append(m.injuryTime, result)
	// }
}

func (m *Match) GetTeams() []*MatchTeam {
	return []*MatchTeam{m.HomeTeam, m.AwayTeam}
}

func (mt *Match) GetMinute() int {
	mt.mut.RLock()
	defer mt.mut.RUnlock()

	return mt.minute
}

func (mt *Match) SetMinute(minute int) {
	mt.mut.RLock()
	defer mt.mut.RUnlock()

	mt.minute = minute
}

func (s *Match) notify(event string) {
	if s.observers != nil {
		zap.L().Debug("notifying observers", zap.String("event", event))
		utils.EachFunc(s.observers, func(o Observer) {
			o.Update(event, s)
		})
	}
}

// func (mt *Match) Subscribe(o Observer) {
// 	mt.mut.Lock()
// 	defer mt.mut.Unlock()

// 	observers := mt.observers
// 	if observers == nil {
// 		observers = []Observer{}
// 	}

// 	if !slices.Contains(observers, o) {
// 		zap.L().Debug("subscribing observer", zap.Any("observer", o))
// 		mt.observers = append(observers, o)
// 	} else {
// 		zap.L().DPanic("observer already subscribed", zap.Any("observer", o))
// 	}
// }

// func (mt *Match) Unsubscribe(observer Observer) {
// 	mt.mut.Lock()
// 	defer mt.mut.Unlock()

// 	currentObservers := mt.observers
// 	if currentObservers != nil {
// 		curLen := len(currentObservers)

// 		currentObservers = slices.DeleteFunc(currentObservers, func(o Observer) bool {
// 			return o == observer
// 		})

// 		if curLen == len(currentObservers) {
// 			zap.L().DPanic("observer not found", zap.Any("observer", observer))
// 		}

// 		zap.L().Debug("unsubscribing observer", zap.Any("observer", observer))
// 		mt.observers = currentObservers
// 	} else {
// 		zap.L().DPanic("no observers found")
// 	}
// }

// func (mt *Match) UnsubscribeAll() {
// 	mt.mut.Lock()
// 	defer mt.mut.Unlock()

// 	mt.observers = []Observer{}
// }

func NewMatch(homeConfig, awayConfig *TeamConfig) *Match {
	return &Match{
		mut:         &sync.RWMutex{},
		minute:      0,
		HomeTeam:    NewMatchTeam(homeConfig),
		AwayTeam:    NewMatchTeam(awayConfig),
		SubjectImpl: NewSubject(),
	}
}

func NewDefaultReferee() *Referee {
	return &Referee{
		Name: defaults.DefaultRefereeName,
		Nat:  defaults.DefaultRefereeNat,
	}
}
