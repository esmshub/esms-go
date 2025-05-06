package formulas

import (
	"math"
)

var (
	ShotOnTargetProbability   = 5000.0 // 50%
	CornerFromSaveProbability = 975.0  // 9.75%
)

func CalcShotProbability(aggression, shooting, passing, oppTackling int) float64 {
	/* Calculate shot probability */
	return 1.65 * (float64(aggression)/50.0 + 800.0*
		math.Pow(((1.0/3.1*float64(shooting)+2.0/3.1*float64(passing))/(float64(oppTackling)+1.0)), 2))
}

func GetShotOnTargetProbability(shooting int) float64 {
	return ShotOnTargetProbability + (float64(shooting) * 50.0)
}
