package models

import "github.com/esmshub/esms-go/engine/events"

type MatchEvent struct {
	event      events.Event
	match      *Match
	activeTeam *MatchTeam
}

func NewMatchEvent(event events.Event, match *Match, activeTeam *MatchTeam) *MatchEvent {
	return &MatchEvent{
		event:      event,
		match:      match,
		activeTeam: activeTeam,
	}
}

func (e *MatchEvent) GetName() string {
	return e.event.GetName()
}

func (e *MatchEvent) GetData() any {
	return e.event.GetData()
}

func (e *MatchEvent) GetMatch() *Match {
	return e.match
}

func (e *MatchEvent) GetActiveTeam() *MatchTeam {
	return e.activeTeam
}

func (e *MatchEvent) GetOppositionTeam() *MatchTeam {
	if e.activeTeam == e.match.HomeTeam {
		return e.match.AwayTeam
	}
	return e.match.HomeTeam
}
