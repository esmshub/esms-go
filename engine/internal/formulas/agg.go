package formulas

func CalcAggression(teamAgg, playerAgg float64) float64 {
	return playerAgg*0.75 + 0.025*teamAgg*playerAgg
}
