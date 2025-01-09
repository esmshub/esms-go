package engine

import "fmt"

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

func Run() {
	// homeTeam := Team{}
	// awayTeam := Team{}

	// load config
	// load teamsheets
	// load rosters
	// load tactics
	// init team / player data

	var gameMinute *int = new(int)
	*gameMinute = 0

	fmt.Println("---------- Kick off ----------")
	runHalf(func(i int) {
		*gameMinute++
		fmt.Println("Minute elapsed:", *gameMinute)
	})
	fmt.Println("---------- Half time ----------")
	runHalf(func(i int) {
		*gameMinute++
		fmt.Println("Minute elapsed:", *gameMinute)
	})
	fmt.Println("---------- Full time ----------")
}
