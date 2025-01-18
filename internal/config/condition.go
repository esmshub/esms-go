package config

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/esmshub/esms-go/engine/types"
	"github.com/esmshub/esms-go/engine/validators"
)

type Condition interface {
	Validate() error
}

type ConditionParser func(args ...string) (Condition, error)

type MinutesEventCondition struct {
	Operator string
	Value    int
}

func (c MinutesEventCondition) Validate() error {
	return nil
}

type InjuryEventCondition struct {
	Value     string
	ValueType string
}

func (c InjuryEventCondition) Validate() error {
	// // TODO: search the roster for player number
	// if util.IsNumericalStr(c.NumberOrPosition) {

	// 	// check if valid squad number
	// } else if isValidPosition(c.NumberOrPosition) {
	// 	// check if valid position
	// } else {
	// 	return fmt.Errorf("invalid position '%s' in %s conditional", c.NumberOrPosition, strings.ToUpper(InjuryEventCode))
	// }
	return nil
}

type CardEventCondition struct {
	CardType  string
	Value     string
	ValueType string
}

func (c CardEventCondition) Validate() error {
	return nil
}

type ShotEventCondition struct {
	Operator string
	Diff     int
}

func (c ShotEventCondition) Validate() error {
	return nil
}

type ScoreEventCondition struct {
	Operator string
	Diff     int
}

func (c ScoreEventCondition) Validate() error {
	return nil
}

func parseInjuryCondition(args ...string) (Condition, error) {
	var c InjuryEventCondition
	if len(args) != 3 {
		return c, fmt.Errorf("invalid %s condition", args[0])
	}
	c.Value = args[2]
	c.ValueType = "position"
	index, err := strconv.Atoi(args[2])
	if err == nil {
		c.ValueType = "number"
		if index < 1 || index > 11 {
			return c, fmt.Errorf("value '%d' is invalid for %s condition, must be between 1-11", index, args[0])
		}
	} else {
		pos := strings.TrimPrefix(c.Value, "O")
		if !validators.IsValidPosition(pos) {
			return c, fmt.Errorf("value '%s' is invalid for %s condition, must be one of %+v", c.Value, args[0], types.ValidPositions)
		}
	}
	return c, nil
}

func parseShotCondition(args ...string) (Condition, error) {
	var c ShotEventCondition
	if len(args) != 3 {
		return c, fmt.Errorf("invalid %s condition", args[0])
	}
	diff, err := strconv.Atoi(args[2])
	if err != nil {
		return c, fmt.Errorf("invalid value '%s' for %s condition, must be valid number", args[2], args[0])
	} else if !slices.ContainsFunc(supportedOperators, func(s string) bool {
		return strings.EqualFold(s, args[1])
	}) {
		return c, fmt.Errorf("invalid operator '%s' for %s condition, must be one of %+v", args[1], args[2], supportedOperators)
	}
	c = ShotEventCondition{
		Operator: args[1],
		Diff:     diff,
	}
	return c, nil
}

func parseScoreCondition(args ...string) (Condition, error) {
	var c ScoreEventCondition
	if len(args) != 3 {
		return c, fmt.Errorf("invalid %s condition", args[0])
	}
	diff, err := strconv.Atoi(args[2])
	if err != nil {
		return c, fmt.Errorf("invalid value '%s' for %s condition, must be valid number", args[2], args[0])
	} else if !slices.ContainsFunc(supportedOperators, func(s string) bool {
		return strings.EqualFold(s, args[1])
	}) {
		return c, fmt.Errorf("invalid operator '%s' for %s condition, must be one of %+v", args[1], args[2], supportedOperators)
	}
	c = ScoreEventCondition{
		Operator: args[1],
		Diff:     diff,
	}
	return c, nil
}

func parseMinutesCondition(args ...string) (Condition, error) {
	var c MinutesEventCondition
	if len(args) != 3 {
		return c, fmt.Errorf("invalid %s condition", args[0])
	}
	min, err := strconv.Atoi(args[2])
	if err != nil {
		return c, fmt.Errorf("invalid value '%s' for %s condition, must be valid number", args[2], args[0])
	} else if !slices.ContainsFunc(supportedOperators, func(s string) bool {
		return strings.EqualFold(s, args[1])
	}) {
		return c, fmt.Errorf("invalid operator '%s' for %s condition, must be one of %+v", args[1], args[2], supportedOperators)
	}
	c = MinutesEventCondition{
		Operator: args[1],
		Value:    min,
	}
	return c, nil
}

func parseCardCondition(args ...string) (Condition, error) {
	var c CardEventCondition
	if len(args) != 3 {
		return c, fmt.Errorf("invalid %s condition", args[0])
	}
	c.CardType = args[0]
	c.Value = args[2]
	c.ValueType = "position"
	pos, err := strconv.Atoi(args[2])
	if err == nil {
		c.ValueType = "number"
		if pos < 1 || pos > 11 {
			return c, fmt.Errorf("value '%d' is invalid for %s condition, must be between 1-11", pos, args[0])
		}
	} else if !validators.IsValidPosition(args[2]) {
		return c, fmt.Errorf("value '%s' is invalid for %s condition, must be one of %+v", args[2], args[0], types.ValidPositions)
	}
	return c, nil
}
