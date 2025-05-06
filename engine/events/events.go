package events

import (
	"reflect"
	"time"
)

const (
	MinuteElapsedEventName               = "MinuteElapsed"
	InjuryTimeAddedEventName             = "InjuryTimeAdded"
	KickOffEventName                     = "KickOff"
	HalfTimeEventName                    = "HalfTime"
	FullTimeEventName                    = "FullTime"
	ShotTackledEventName                 = "ShotTackled"
	ShotTackledCornerEventName           = "ShotTackledCorner"
	ShotTackledRecoveryEventName         = "ShotTackledRecovery"
	ShotSavedEventName                   = "ShotSaved"
	ShotSavedCornerEventName             = "ShotSavedCorner"
	ShotOffTargetEventName               = "ShotOffTarget"
	ShotOffTargetDeflectionEventName     = "ShotOffTargetDeflection"
	ShotClearedEventName                 = "ShotCleared"
	ShotOnTargetEventName                = "ShotOnTarget"
	ChanceEventName                      = "Chance"
	ChanceBeatsDefenderEventName         = "ChanceBeatsDefender"
	AssistedChanceEventName              = "AssistedChance"
	AssistedChanceBeatsDefenderEventName = "AssistedChanceBeatsDefender"
	GoalScoredEventName                  = "GoalScored"
	GoalScoredCancelledEventName         = "GoalScoredCancelled"
	OwnGoalScoredEventName               = "OwnGoalScored"
	CornerCaughtEventName                = "CornerCaught"
	CornerClearedEventName               = "CornerCleared"
	CornerShotEventName                  = "CornerShot"
)

type Event interface {
	GetName() string
	GetData() any
}

// Event represents a match event
type EventImpl struct {
	data      any
	name      string
	Timestamp time.Time
}

func clone(v any) any {
	if v == nil {
		return nil
	}

	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	switch typ.Kind() {
	case reflect.Ptr:
		// Clone pointer to value
		elem := val.Elem()
		if !elem.IsValid() {
			return v // nil pointer
		}
		clone := reflect.New(elem.Type())
		clone.Elem().Set(elem)
		return clone.Interface()

	case reflect.Slice:
		// Clone slice
		clone := reflect.MakeSlice(typ, val.Len(), val.Cap())
		reflect.Copy(clone, val)
		return clone.Interface()

	case reflect.Map:
		// Clone map
		clone := reflect.MakeMapWithSize(typ, val.Len())
		for _, key := range val.MapKeys() {
			clone.SetMapIndex(key, val.MapIndex(key))
		}
		return clone.Interface()

	case reflect.Chan, reflect.Func:
		// Cannot deep-clone channels or functions
		return v

	default:
		// Value type (int, struct, string, etc.)
		return v
	}
}

func NewEvent(name string, data any) Event {
	return &EventImpl{
		data:      clone(data),
		name:      name,
		Timestamp: time.Now(),
	}
}

func (e *EventImpl) GetData() any {
	return e.data
}

func (e *EventImpl) GetName() string {
	return e.name
}

// EventHandler defines how to handle an event
type EventHandler interface {
	Handle(event Event) ([]Event, error)
}
