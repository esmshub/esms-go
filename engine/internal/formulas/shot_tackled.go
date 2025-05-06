package formulas

var ShotTackledBaseProbability = 3200 // 32%

func CalcShotTackledProbability(shooting, passing, oppTackling int) float64 {
	return float64(ShotTackledBaseProbability) * ((float64(oppTackling) * 3.0) / (float64(passing)*2.0 + float64(shooting)))
}
