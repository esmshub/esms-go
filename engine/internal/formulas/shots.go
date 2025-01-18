package formulas

import (
	"math"
)

func CalcShotProbability(aggression, shooting, passing, defensiveTackling float64) float64 {
	/* Calculate shot probability */
	return 1.65 * (float64(aggression)/50.0 + 800.0*
		math.Pow(((1.0/3.1*float64(shooting)+2.0/3.1*float64(passing))/(float64(defensiveTackling)+1.0)), 2))
}
