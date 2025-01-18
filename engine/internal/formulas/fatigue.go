package formulas

const FATIGUE_DEGRADATION_RATE = 0.996

func CalcFatigue(current float64, minute int) float64 {
	return current * FATIGUE_DEGRADATION_RATE
	// return current * math.Pow(FATIGUE_DEGRADATION_RATE, factor)
}
