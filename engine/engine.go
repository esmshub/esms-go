package engine

import (
	"fmt"

	"github.com/esmshub/esms-go/engine/models"
	"github.com/esmshub/esms-go/engine/pkg/rng"
	"github.com/esmshub/esms-go/engine/types"
	"github.com/esmshub/esms-go/engine/validators"
	"go.uber.org/zap"
)

type MinuteElapsedHandler func(int)

func updateTeamAttrs(config *models.TeamConfig, tactics *models.TacticsMatrix) {
	config.TeamAbility.Goalkeeping = 0
	config.TeamAbility.Tackling = 0
	config.TeamAbility.Passing = 0
	config.TeamAbility.Shooting = 0

	players := append(config.Lineup, config.Subs...)
	for _, p := range players {
		if p.IsActive {
			if p.Position == types.POSITION_GK {
				p.Ability.Goalkeeping = p.BaseAbility.Goalkeeping
				config.TeamAbility.Goalkeeping += p.Ability.Goalkeeping
			} else if tactics != nil {
				matrix := (*tactics)[fmt.Sprintf("%s_%s", config.Tactic, p.Position)]
				p.Ability.Tackling = int((matrix[0] * float64(p.BaseAbility.Tackling)) * p.Fatigue)
				p.Ability.Passing = int((matrix[1] * float64(p.BaseAbility.Passing)) * p.Fatigue)
				p.Ability.Shooting = int((matrix[2] * float64(p.BaseAbility.Shooting)) * p.Fatigue)

				// update cumulative stats
				config.TeamAbility.Tackling += p.Ability.Tackling
				config.TeamAbility.Passing += p.Ability.Passing
				config.TeamAbility.Shooting += p.Ability.Shooting
			}
		} else {
			p.Ability.Goalkeeping = 0
			p.Ability.Tackling = 0
			p.Ability.Passing = 0
			p.Ability.Shooting = 0
		}
	}
}

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

func Run(match *models.Match, options *Options) (*models.MatchResult, error) {
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

	validator := validators.NewTeamConfigValidator(options.AppConfig)
	if err := validator.Validate(match.HomeTeam); err != nil {
		return nil, err
	}
	if err := validator.Validate(match.AwayTeam); err != nil {
		return nil, err
	}

	// initialize team attributes
	updateTeamAttrs(match.HomeTeam, options.TacticsMatrix)
	updateTeamAttrs(match.AwayTeam, options.TacticsMatrix)

	// errors := match.HomeTeam.Validate()
	// errors = append(errors, match.AwayTeam.Validate()...)
	// zap.L().Info("Running fixture", zap.Any("fixture", fixture), zap.Any("options", options))
	// validate fixtureset

	var gameMinute *int = new(int)
	*gameMinute = 0

	fmt.Println("---------- Kick off ----------")
	fmt.Println("Home bonus:", options.AppConfig["home_bonus"])
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
	}, nil
}
