package formulas

func GetOneOnOneProbability(factor int) float64 {
	return 250.0 + (float64(factor) * 550.0)
}
