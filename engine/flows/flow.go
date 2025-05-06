package flows

import (
	"github.com/esmshub/esms-go/engine/events"
	"github.com/esmshub/esms-go/engine/models"
	"github.com/esmshub/esms-go/engine/pkg/rng"
)

type Flow func(*models.MatchTeam, *models.MatchTeam) []events.Event

var All = [3]Flow{ChanceFlow, FoulFlow, InjuryFlow}

func getAttacker(team *models.MatchTeam) *models.MatchPlayer {
	return getWeightedPlayer(team.GetActive(), func(p *models.MatchPlayer) int {
		return (p.GetMatchAbility().Shooting * 25) + (p.GetBaseAbility().Shooting * 10)
	})
}

func getAssister(team *models.MatchTeam) *models.MatchPlayer {
	return getWeightedPlayer(team.GetActive(), func(p *models.MatchPlayer) int {
		return (p.GetMatchAbility().Passing * 240) + (p.GetBaseAbility().Passing * 100)
	})
}

func getDefender(team *models.MatchTeam) *models.MatchPlayer {
	return getWeightedPlayer(team.GetActive(), func(p *models.MatchPlayer) int {
		return (p.GetMatchAbility().Tackling * 25) + (p.GetBaseAbility().Tackling * 10)
	})
}

func getWeightedPlayer(players []*models.MatchPlayer, weightFunc func(*models.MatchPlayer) int) *models.MatchPlayer {
	// Calculate total weight
	totalWeight := 0
	weights := []int{}
	for _, p := range players {
		val := weightFunc(p)
		weights = append(weights, val)
		totalWeight += val
	}

	// Generate a random value between 0 and total weight
	threshold := rng.Random(totalWeight)

	// Determine the player based on the random value and cumulative weight
	cumulativeWeight := 0
	for _, p := range players {
		cumulativeWeight += weightFunc(p)
		if cumulativeWeight > threshold {
			return p
		}
	}

	panic("No player selected")
}
