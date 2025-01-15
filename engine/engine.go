package engine

import (
	"fmt"

	"github.com/esmshub/esms-go/engine/models"
	"github.com/esmshub/esms-go/engine/pkg/rng"
	"go.uber.org/zap"
)

type MinuteElapsedHandler func(int)

func runHalf(minuteElapsedHandler MinuteElapsedHandler) {
	for i := 0; i < 45; i++ {
		// - did a shot happen?
		// -- was it off-target?
		// --- did it hit the post?
		// --- did it hit the bar?
		// --- did it go over the bar?
		// --- did it go past the post?
		// -- was it on-target?
		// --- was it blocked?
		// ---- was it a corner?
		// --- was it saved?
		// ---- was it collected?
		// ---- was it a corner?
		// ---- was it parried?
		// --- was it a goal?
		// ---- who scored?
		// ---- was there an assist?
		// ---- was it an own goal?
		// ---- is there a VAR review?
		// ----- was it offside?
		// ----- was it a handball?
		// ----- was it a foul?

		// - did a tackle happen?
		// -- was it a foul?
		// -- was it a yellow card?v
		if minuteElapsedHandler != nil {
			minuteElapsedHandler(i)
		}
	}
}

func Run(match *models.Match, options *Options) *models.MatchResult {
	if options.RngSeed != 0 {
		// seed random number generator
		zap.L().Info("Seeding RNG", zap.Uint64("seed", options.RngSeed))
		rng.Seed(options.RngSeed)
	}
	if match.Referee == nil {
		zap.L().Info("No referee assigned, using default")
		match.Referee = &models.Referee{
			Name: "Perluigi Collina",
			Nat:  "Italian",
		}
	}
	// errors := match.HomeTeam.Validate()
	// errors = append(errors, match.AwayTeam.Validate()...)
	// zap.L().Info("Running fixture", zap.Any("fixture", fixture), zap.Any("options", options))
	// validate fixtureset

	// homeTeam := Team{}
	// awayTeam := Team{}

	// load teamsheets
	// load rosters
	// load tactics
	// init team / player data

	var gameMinute *int = new(int)
	*gameMinute = 0

	fmt.Println("---------- Kick off ----------")
	fmt.Println("Home bonus:", options.HomeBonus)
	fmt.Println("Match type:", options.MatchType)
	// runHalf(func(i int) {
	// 	*gameMinute++
	// 	fmt.Println("Minute elapsed:", *gameMinute)
	// })
	// fmt.Println("---------- Half time ----------")
	// runHalf(func(i int) {
	// 	*gameMinute++
	// 	fmt.Println("Minute elapsed:", *gameMinute)
	// })
	// fmt.Println("---------- Full time ----------")
	return &models.MatchResult{
		HomeTeam: match.HomeTeam,
		AwayTeam: match.AwayTeam,
		Referee:  match.Referee,
		RngSeed:  rng.GetSeed(),
	}
}
