package config

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/esmshub/esms-go/engine/models"
	"github.com/esmshub/esms-go/engine/types"
	"github.com/esmshub/esms-go/engine/validators"
	"github.com/esmshub/esms-go/pkg/utils"
	"golang.org/x/exp/slices"
)

// func (c Condition) Validate() error {
// 	_, validEvent := util.FindInSlice(types.ValidConditionalEvents, func(s string) bool {
// 		return strings.EqualFold(s, c.Event)
// 	})
// 	if !validEvent {
// 		return fmt.Errorf("invalid event type: %s", c.Event)
// 	}

// 	_, validOp := util.FindInSlice(supportedOperators, func(s string) bool {
// 		return strings.EqualFold(s, c.Operator)
// 	})
// 	if !validOp {
// 		return fmt.Errorf("invalid operator type: %s", c.Operator)
// 	}
// 	if _, err := strconv.Atoi(c.Value); err != nil {
// 		return fmt.Errorf("%s condition values must be numeric", c.Event)
// 	}

// 	return nil
// }

// func (c EventCondition) Evaluate(input any) (bool, error) {
// 	if _, ok := input.(int); !ok {
// 		return false, fmt.Errorf("%s condition values must be numeric", c.Event)
// 	}

// 	value, err := strconv.Atoi(c.Value)
// 	if err != nil {
// 		return false, fmt.Errorf("%s condition values must be numeric", c.Event)
// 	}

// 	switch c.Operator {
// 	case "<":
// 		return input.(int) < value, nil
// 	case "<=":
// 		return input.(int) <= value, nil
// 	case ">=":
// 		return input.(int) >= value, nil
// 	case ">":
// 		return input.(int) > value, nil
// 	default:
// 		return input.(int) == value, nil
// 	}
// }

type ConditionalParser func(string) (*models.Conditional, error)

func parseAggression(text string, action string) (int, error) {
	parts := strings.Fields(text)
	if len(parts) < 2 {
		return -1, fmt.Errorf("invalid %s conditional: %s", action, text)
	}
	agg, err := strconv.Atoi(parts[1])
	if err != nil {
		return -1, fmt.Errorf("%s value '%s' is not a valid number", action, parts[1])
	} else if agg < 1 || agg > 20 {
		return -1, fmt.Errorf("%s value must be between 1-20", action)
	}

	return agg, nil
}

func parseAggressionConditional(text string) (*models.Conditional, error) {
	var cond models.Conditional
	value, err := parseAggression(text, types.AggressionAction)
	if err != nil {
		return &cond, err
	} else {
		cond = models.Conditional{
			Action: types.AggressionAction,
			Values: []any{value},
		}
		return &cond, nil
	}
}

func parseChangeAggressionConditional(text string) (*models.Conditional, error) {
	var cond *models.Conditional
	parts := strings.Split(text, "IF")
	if len(parts) != 2 {
		return cond, fmt.Errorf("invalid %s conditional: %s", types.ChangeAggressionAction, text)
	}
	agg, err := parseAggression(parts[0], types.ChangeAggressionAction)
	if err != nil {
		return cond, err
	}

	conditions, err := parseConditions(parts[1])
	if err != nil {
		return cond, err
	} else if len(conditions) == 0 {
		return cond, fmt.Errorf("no valid conditions found for %s conditional '%s'", types.ChangeAggressionAction, text)
	}

	cond = &models.Conditional{
		Action:     types.ChangeAggressionAction,
		Values:     []any{agg},
		Conditions: conditions,
	}
	return cond, nil
}

func parseChangePositionConditional(text string) (*models.Conditional, error) {
	var cond *models.Conditional
	parts := strings.Split(text, "IF")
	if len(parts) != 2 {
		return cond, fmt.Errorf("invalid %s conditional: %s", types.ChangePositionAction, text)
	}

	fields := strings.Fields(parts[0])
	if len(fields) < 3 {
		return cond, fmt.Errorf("invalid %s conditional: %s", types.ChangePositionAction, text)
	}
	pos := strings.TrimSpace(fields[2])
	num, err := strconv.Atoi(strings.TrimSpace(fields[1]))
	if err != nil || (num < 0 || num > 11) {
		return cond, fmt.Errorf("invalid value '%s' in %s conditional, must be a number between 1-11", strings.TrimSpace(fields[1]), types.ChangePositionAction)
	} else if !validators.IsValidPosition(pos) {
		return cond, fmt.Errorf("invalid value '%s' in %s conditional, must be one of %+v", pos, types.ChangePositionAction, types.ValidPositions)
	}

	conditions, err := parseConditions(text)
	if err != nil {
		return cond, err
	} else if len(conditions) == 0 {
		return cond, fmt.Errorf("no valid conditions found for %s conditional '%s'", types.ChangePositionAction, text)
	}

	cond = &models.Conditional{
		Action:     types.ChangePositionAction,
		Values:     []any{num, pos},
		Conditions: conditions,
	}
	return cond, nil
}

func parsePlayerNumberOrPosition(text string) (any, error) {
	str := strings.TrimSpace(text)
	num, err := strconv.Atoi(str)
	if err == nil {
		return num, nil
	} else if !validators.IsValidPosition(str) {
		return str, fmt.Errorf("must be one of %+v", types.ValidPositions)
	} else {
		return str, nil
	}
}

func parseSubConditional(text string) (*models.Conditional, error) {
	var cond *models.Conditional
	parts := strings.Split(text, "IF")
	if len(parts) != 2 {
		return cond, fmt.Errorf("invalid %s conditional: %s", types.SubstituteAction, text)
	}

	fields := strings.Fields(parts[0])
	if len(fields) < 4 {
		return cond, fmt.Errorf("invalid %s conditional: %s", types.SubstituteAction, text)
	}
	activeNumOrPos, err := parsePlayerNumberOrPosition(fields[1])
	if err != nil {
		return cond, fmt.Errorf("value '%+v' invalid for %s condition, %+v", activeNumOrPos, types.SubstituteAction, err)
	} else if utils.IsNumber(activeNumOrPos) && (activeNumOrPos.(int) < 1 || activeNumOrPos.(int) > 11) {
		return cond, fmt.Errorf("value '%+v' invalid for %s condition, must be a number between 1-11", activeNumOrPos, types.SubstituteAction)
	}
	subNumOrPos, err := strconv.Atoi(strings.TrimSpace(fields[2]))
	if err != nil || (subNumOrPos < 12 || subNumOrPos > 16) {
		return cond, fmt.Errorf("value '%s' invalid for %s condition, must be a number between 12-16", strings.TrimSpace(fields[2]), types.SubstituteAction)
	}
	// parse target position
	targetPos := strings.TrimSpace(fields[3])
	if !validators.IsValidPosition(targetPos) {
		return cond, fmt.Errorf("value '%s' invalid for %s conditional, must be one of %+v", targetPos, types.SubstituteAction, types.ValidPositions)
	}

	conditions, err := parseConditions(parts[1])
	if err != nil {
		return cond, err
	} else if len(conditions) == 0 {
		return cond, fmt.Errorf("no conditions found for %s conditional '%s'", types.SubstituteAction, text)
	}

	cond = &models.Conditional{
		Action:     types.SubstituteAction,
		Values:     []any{activeNumOrPos, subNumOrPos, targetPos},
		Conditions: conditions,
	}
	return cond, nil
}

func parseTacticConditional(text string) (*models.Conditional, error) {
	var cond *models.Conditional
	parts := strings.Split(text, "IF")
	if len(parts) != 2 {
		return cond, fmt.Errorf("invalid %s conditional: %s", types.ChangeTacticAction, text)
	}

	fields := strings.Fields(parts[0])
	if len(fields) < 2 {
		return cond, fmt.Errorf("invalid %s conditional: %s", types.ChangeTacticAction, text)
	}

	tactic := strings.TrimSpace(fields[1])
	if !isValidTactic(tactic) {
		return cond, fmt.Errorf("value '%s' is invalid for %s conditional, must be one of %+v", tactic, types.ChangeTacticAction, types.ValidTactics)
	}

	conditions, err := parseConditions(parts[1])
	if err != nil {
		return cond, err
	} else if len(conditions) == 0 {
		return cond, fmt.Errorf("no conditions found for %s conditional '%s'", types.ChangeTacticAction, text)
	}

	cond = &models.Conditional{
		Action:     types.ChangeTacticAction,
		Values:     []any{tactic},
		Conditions: conditions,
	}
	return cond, nil
}

func parseConditions(text string) ([]models.Condition, error) {
	results := []models.Condition{}
	matches := conditionRegex.FindAllStringSubmatch(strings.TrimSpace(text), -1)
	for _, match := range matches {
		idx := slices.IndexFunc(types.ValidConditionalEvents, func(s string) bool {
			return strings.EqualFold(s, match[1])
		})
		if idx == -1 {
			return results, fmt.Errorf("invalid event type: %s", match[1])
		}

		parse := conditionParsers[types.ValidConditionalEvents[idx]]
		if parse == nil {
			panic(fmt.Errorf("no parser for %s condition found", types.ValidConditionalEvents[idx]))
		}
		cond, err := parse(match[1:]...)
		if err != nil {
			return results, err
		}
		// if err := cond.Validate(); err != nil {
		// 	return results, err
		// }
		results = append(results, cond)
	}
	return results, nil
}

var (
	supportedOperators    []string = []string{"<", "<=", "=", ">=", ">"}
	conditionRegexPattern string   = "\\b(YELLOW|RED|SHOTS|SCORE|MIN|INJURED)\\b(?:\\s+)?(?:([<>]?=?)+\\s*(-?\\d*\\.?\\d+|-?[a-zA-Z]+))\\b"
	// conditionalRegexPattern string   = "\\b(AGG|CHANGEAGG|SUB|TACTIC|)\\b(?:\\s+)(.*)(?=IF)"
	// conditionRegexPattern string   = "\\b(YELLOW|RED|INJURED|SCORE|SHOTS|MIN)\\b\\s*([^\\s]+?)?\\s*([-]?[a-zA-Z]+(?:[0-9]*\\.[0-9]+)?|[-]?[0-9]+(?:\\.[0-9]+)?)"
	// regex string = "(?:.)(YELLOW|RED|INJURED|SCORE|SHOTS|MIN)(?:\s+)([<>=]+)?(?:\s+)?([a-zA-Z]+|[0-9]+)(?:\s+)?"
	// conditionalRegex *regexp.Regexp
	conditionRegex *regexp.Regexp

	conditionParsers = map[string]ConditionParser{
		types.InjuryEvent:     parseInjuryCondition,
		types.ShotEvent:       parseShotCondition,
		types.ScoreEvent:      parseScoreCondition,
		types.MinuteEvent:     parseMinutesCondition,
		types.RedCardEvent:    parseCardCondition,
		types.YellowCardEvent: parseCardCondition,
	}

	conditionalParsers = map[string]ConditionalParser{
		types.AggressionAction:       parseAggressionConditional,
		types.ChangeAggressionAction: parseChangeAggressionConditional,
		types.ChangePositionAction:   parseChangePositionConditional,
		types.SubstituteAction:       parseSubConditional,
		types.ChangeTacticAction:     parseTacticConditional,
	}
)

func init() {
	// conditionalRegex = regexp.MustCompile(conditionalRegexPattern)
	conditionRegex = regexp.MustCompile(conditionRegexPattern)
}
