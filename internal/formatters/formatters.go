package formatters

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/esmshub/esms-go/engine/models"
	"github.com/esmshub/esms-go/pkg/utils"
	"golang.org/x/exp/maps"
)

type FormatScorersOptions struct {
	RowDelimiter         string
	NoResultsPlaceholder string
}

func FormatScorers(scorers []models.MatchGoal, options FormatScorersOptions) string {
	if len(scorers) == 0 {
		return options.NoResultsPlaceholder
	}

	// Group scorers by name and minute
	scorersByMinute := utils.Reduce(scorers, func(acc map[string][]string, g models.MatchGoal) map[string][]string {
		name := g.Player.GetName()
		if _, ok := acc[name]; !ok {
			acc[name] = []string{}
		}
		minStr := strconv.Itoa(g.Minute)
		if g.IsOwnGoal {
			minStr = minStr + " og"
		} else if g.IsPenalty {
			minStr = minStr + " pen"
		}
		acc[name] = append(acc[name], minStr)
		return acc
	}, map[string][]string{})

	return strings.Join(utils.Map(maps.Keys(scorersByMinute), func(name string) string {
		return fmt.Sprintf("%s (%s)", name, strings.Join(scorersByMinute[name], ","))
	}), options.RowDelimiter)
}
