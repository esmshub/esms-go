package formulas

var GoalProbability = 2200.0       // 22%
var GoalCancelledProbability = 350 // 3.5%

func GetGoalProbability(shooting, oppKeeping int) float64 {
	attacking := float64(shooting) * 240.0
	defense := float64(oppKeeping) * 170.0
	return attacking - defense + GoalProbability
}
