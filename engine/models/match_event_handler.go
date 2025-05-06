package models

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/esmshub/esms-go/engine/events"
)

type MatchEventHandler struct{}

func getData[T any](event events.Event) T {
	data, ok := event.GetData().(T)
	if !ok {
		var t T
		panic(fmt.Sprintf("data for %s event is %s, expected %s", event.GetName(), reflect.TypeOf(data).Name(), reflect.TypeOf(t).Name()))
	}
	return data
}

func getElement[T any](v map[string]any, key string) T {
	element, ok := v[key]
	if !ok {
		panic(fmt.Sprintf("element '%s' does not exist in map", key))
	}

	var result T
	result, ok = element.(T)
	if !ok {
		panic(fmt.Sprintf("element '%s' is not of type %s", key, reflect.TypeOf(element).Name()))
	}

	return result
}

func (h MatchEventHandler) Handle(event events.Event) ([]events.Event, error) {
	matchEvent, ok := event.(*MatchEvent)
	if !ok {
		return []events.Event{}, errors.New("event is not a MatchEvent")
	}

	switch event.GetName() {
	case events.OwnGoalScoredEventName:
		return h.handleOwnGoalScored(matchEvent)
	case events.AssistedChanceEventName, events.AssistedChanceBeatsDefenderEventName:
		return h.HandleAssistedChance(matchEvent)
	case events.ShotTackledEventName, events.ShotTackledCornerEventName, events.ShotTackledRecoveryEventName:
		return h.handleShotTackled(matchEvent)
	case events.ShotOnTargetEventName:
		return h.handleShotOnTarget(matchEvent)
	case events.ShotOffTargetEventName, events.ShotOffTargetDeflectionEventName:
		return h.handleShotOffTarget(matchEvent)
	case events.ShotSavedEventName, events.ShotSavedCornerEventName, events.ShotClearedEventName:
		return h.handleShotSaved(matchEvent)
	case events.GoalScoredEventName:
		return h.handleGoalScored(matchEvent)
	}

	return []events.Event{}, nil
}

func (h *MatchEventHandler) HandleAssistedChance(matchEvent *MatchEvent) ([]events.Event, error) {
	data := getData[map[string]any](matchEvent)
	assister := getElement[*MatchPlayer](data, "assister")
	assister.AddKeyPass()
	matchEvent.GetActiveTeam().AddPass()

	return []events.Event{}, nil
}

func (h *MatchEventHandler) handleOwnGoalScored(matchEvent *MatchEvent) ([]events.Event, error) {
	match := matchEvent.GetMatch()

	minute := match.GetMinute()
	data := getData[map[string]any](matchEvent)
	scorer := getElement[*MatchPlayer](data, "scorer")
	scorer.AddOwnGoal(minute)
	goalkeeper := getElement[*MatchPlayer](data, "opp_keeper")
	goalkeeper.AddConceded()
	// update match stats
	matchEvent.GetActiveTeam().AddGoal(scorer, minute)

	return []events.Event{}, nil
}

func (h *MatchEventHandler) handleShotTackled(matchEvent *MatchEvent) ([]events.Event, error) {
	data := getData[map[string]any](matchEvent)
	defender := getElement[*MatchPlayer](data, "tackler")
	defender.AddKeyTackle()
	// update match stats
	match := matchEvent.GetMatch()
	if matchEvent.GetActiveTeam().GetShortName() == match.HomeTeam.GetShortName() {
		match.AwayTeam.AddTackle()
	} else {
		match.HomeTeam.AddTackle()
	}
	return []events.Event{}, nil
}

func (h *MatchEventHandler) handleShotOnTarget(matchEvent *MatchEvent) ([]events.Event, error) {
	data := getData[map[string]any](matchEvent)
	attacker := getElement[*MatchPlayer](data, "attacker")
	attacker.AddShotOnTarget()
	// update match stats
	matchEvent.GetActiveTeam().AddShotOnTarget()
	return []events.Event{}, nil
}

func (h *MatchEventHandler) handleShotOffTarget(matchEvent *MatchEvent) ([]events.Event, error) {
	data := getData[map[string]any](matchEvent)
	attacker := getElement[*MatchPlayer](data, "attacker")
	attacker.AddShotOffTarget()
	// update match stats
	matchEvent.GetActiveTeam().AddShotOffTarget()
	return []events.Event{}, nil
}

func (h *MatchEventHandler) handleShotSaved(matchEvent *MatchEvent) ([]events.Event, error) {
	data := getData[map[string]any](matchEvent)
	keeper := getElement[*MatchPlayer](data, "opp_keeper")
	keeper.AddSave()
	// update match stats
	matchEvent.GetActiveTeam().AddShotOnTarget()
	return []events.Event{}, nil
}

func (h *MatchEventHandler) handleGoalScored(matchEvent *MatchEvent) ([]events.Event, error) {
	data := getData[map[string]any](matchEvent)
	wasCancelled := getElement[bool](data, "goal_cancelled")
	if !wasCancelled {
		scorer := getElement[*MatchPlayer](data, "attacker")
		minute := matchEvent.GetMatch().GetMinute()
		// update goal stats
		matchEvent.GetActiveTeam().AddGoal(scorer, minute)
		scorer.AddGoal(minute)
		// update conceded stats
		keeper := getElement[*MatchPlayer](data, "opp_keeper")
		keeper.AddConceded()
		// update assist stats
		if assister, ok := data["assister"].(*MatchPlayer); ok {
			assister.AddAssist()
		}
	}

	return []events.Event{}, nil
}
