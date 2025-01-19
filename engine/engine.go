package engine

import (
	"fmt"
	"math"
	"strings"

	"github.com/esmshub/esms-go/engine/models"
	"github.com/esmshub/esms-go/engine/pkg/rng"
	"github.com/esmshub/esms-go/engine/pkg/utils"
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
	config.TeamAbility.Aggression = 0

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
			}
			// update cumulative stats
			config.TeamAbility.Tackling += p.Ability.Tackling
			config.TeamAbility.Passing += p.Ability.Passing
			config.TeamAbility.Shooting += p.Ability.Shooting
			config.TeamAbility.Aggression += p.Ability.Aggression
		} else {
			p.Ability.Goalkeeping = 0
			p.Ability.Tackling = 0
			p.Ability.Passing = 0
			p.Ability.Shooting = 0
		}
	}
	// initialize team attributes
	// zap.L().Info("Team aggression before", zap.Any("value", config.GetAggression()))
	if aggCond := utils.FindFunc(config.Conditionals, func(c *models.Conditional) bool {
		return strings.EqualFold(c.Action, types.AGG_ACTION)
	}); aggCond != nil {
		config.CalcAggression(aggCond.Values[0].(int))
	} else {
		config.CalcAggression(DEFAULT_AGG_LEVEL)
	}
	// zap.L().Info("Team aggression after", zap.Any("value", config.GetAggression()))
}

// Calculates how much injury time to add.
//
// Takes into account substitutions, injuries and fouls (by both teams)
func calcMaxInjuryTime(match *models.Match) int {
	/*
		for consideration:
		- The assessment of apparently injured players
		- The removal from the field of injured players
		- Substitutions
		- Perceived time wasting by players
		- Red or yellow cards being issued
		- Delays for VAR checks
		- Drinks breaks in hotter venues
	*/
	return rng.RandomRange(1, 5)
	// teams := []*models.TeamConfig{match.HomeTeam, match.AwayTeam}
	// subCount := utils.Reduce(teams, func(acc int, t *models.TeamConfig) int {
	// 	return acc + utils.CountFunc(t.Lineup, func(p *models.PlayerPosition) bool {
	// 		return p.Stats.IsSubbed
	// 	}) + utils.CountFunc(t.Subs, func(p *models.PlayerPosition) bool {
	// 		return p.Stats.IsSubbed
	// 	})
	// }, 0)
	// injuryCount := utils.Reduce(teams, func(acc int, t *models.TeamConfig) int {
	// 	return acc + utils.CountFunc(t.Lineup, func(p *models.PlayerPosition) bool {
	// 		return p.Stats.IsInjured
	// 	}) + utils.CountFunc(t.Subs, func(p *models.PlayerPosition) bool {
	// 		return p.Stats.IsInjured
	// 	})
	// }, 0)
	// foulCount := utils.Reduce(teams, func(acc int, t *models.TeamConfig) int {
	// 	return acc + utils.SumFunc(t.Lineup, func(p *models.PlayerPosition) int {
	// 		return p.Stats.Fouls
	// 	}) + utils.SumFunc(t.Subs, func(p *models.PlayerPosition) int {
	// 		return p.Stats.Fouls
	// 	})
	// }, 0)

	// return int(math.Ceil(float64(subCount+injuryCount+foulCount) * 0.5))
}

// Function to calculate possession
func calcHomePossession(homeAbility, awayAbility *models.PlayerAbilities) float64 {

	// Normalize the stats for each team
	normalizedPassingHome := float64(homeAbility.Passing) / MAX_TEAMSKILL_VALUE
	normalizedTacklingHome := float64(homeAbility.Tackling) / MAX_TEAMSKILL_VALUE
	normalizedShootingHome := float64(homeAbility.Shooting) / MAX_TEAMSKILL_VALUE

	normalizedPassingAway := float64(awayAbility.Tackling) / MAX_TEAMSKILL_VALUE
	normalizedTacklingAway := float64(awayAbility.Passing) / MAX_TEAMSKILL_VALUE
	normalizedShootingAway := float64(awayAbility.Shooting) / MAX_TEAMSKILL_VALUE

	// Calculate the initial possession
	totalStrength := normalizedPassingHome + normalizedTacklingHome + normalizedShootingHome +
		normalizedPassingAway + normalizedTacklingAway + normalizedShootingAway

	poss := (normalizedPassingHome + normalizedTacklingHome + normalizedShootingHome) / totalStrength * 100

	// Add a random factor for variability
	randomFactor := (0.05 + rng.RandomF()*0.10) // A random factor between 0.05 and 0.15
	return math.Min(math.Floor(poss*(1+randomFactor)), 85.0)
}

func runHalf(match *models.Match, minuteElapsedHandler MinuteElapsedHandler) models.MatchStats {
	injuryTime := calcMaxInjuryTime(match)
	stats := models.MatchStats{
		// HomeStats: models.TeamStats{},
		// AwayStats: models.TeamStats{},
	}
	for i := 0; i < (MINS_PER_HALF + injuryTime); i++ {
		if i == MINS_PER_HALF-1 {
			fmt.Printf("%d injury time added\n", injuryTime)
		}
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
		// -- was it a yellow card?

		players := []*models.PlayerPosition{}
		players = append(players, match.HomeTeam.Lineup...)
		players = append(players, match.HomeTeam.Subs...)
		players = append(players, match.AwayTeam.Lineup...)
		players = append(players, match.AwayTeam.Subs...)
		utils.EachFunc(players, func(p *models.PlayerPosition) {
			p.Stats.MinutesPlayed++
		})
		if minuteElapsedHandler != nil {
			minuteElapsedHandler(i)
		}
	}
	return stats
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

	updateTeamAttrs(match.HomeTeam, options.TacticsMatrix)
	match.HomeTeam.CalcShotProbability(match.AwayTeam.TeamAbility)
	updateTeamAttrs(match.AwayTeam, options.TacticsMatrix)
	match.AwayTeam.CalcShotProbability(match.HomeTeam.TeamAbility)

	// errors := match.HomeTeam.Validate()
	// errors = append(errors, match.AwayTeam.Validate()...)
	// zap.L().Info("Running fixture", zap.Any("fixture", fixture), zap.Any("options", options))
	// validate fixtureset

	var gameMinute *int = new(int)
	*gameMinute = 0

	fmt.Println("---------- Kick off ----------")
	zap.L().Info("Initial Team ability", zap.Any("home", match.HomeTeam.TeamAbility))
	zap.L().Info("Initial Team ability", zap.Any("away", match.AwayTeam.TeamAbility))
	zap.L().Info("Before probability", zap.Any("shot_prob", match.HomeTeam.GetShotProbability()))
	halfTimeStats := runHalf(match, func(i int) {
		curMin := i + 1
		if curMin <= MINS_PER_HALF {
			*gameMinute++
		}
		fmt.Println("Minute elapsed:", curMin)

		updateTeamAttrs(match.HomeTeam, options.TacticsMatrix)
		match.HomeTeam.CalcShotProbability(match.AwayTeam.TeamAbility)
		updateTeamAttrs(match.AwayTeam, options.TacticsMatrix)
		match.AwayTeam.CalcShotProbability(match.HomeTeam.TeamAbility)
	})
	fmt.Println("---------- Half time ----------")
	zap.L().Info("Half time stats", zap.Any("stats", halfTimeStats))
	zap.L().Info("After Team ability", zap.Any("home", match.HomeTeam.TeamAbility))
	zap.L().Info("After Team ability", zap.Any("away", match.AwayTeam.TeamAbility))
	zap.L().Info("After Shot probability", zap.Any("shot_prob", match.HomeTeam.GetShotProbability()))
	fullTimeStats := runHalf(match, func(i int) {
		*gameMinute++
		fmt.Println("Minute elapsed:", *gameMinute)

		updateTeamAttrs(match.HomeTeam, options.TacticsMatrix)
		match.HomeTeam.CalcShotProbability(match.AwayTeam.TeamAbility)
		updateTeamAttrs(match.AwayTeam, options.TacticsMatrix)
		match.AwayTeam.CalcShotProbability(match.HomeTeam.TeamAbility)
	})
	fmt.Println("---------- Full time ----------")
	zap.L().Info("Half time stats", zap.Any("stats", fullTimeStats))
	zap.L().Info("After Team ability", zap.Any("home", match.HomeTeam.TeamAbility))
	zap.L().Info("After Team ability", zap.Any("away", match.AwayTeam.TeamAbility))
	zap.L().Info("After Shot probability", zap.Any("shot_prob", match.HomeTeam.GetShotProbability()))
	// calc possession
	homePoss := calcHomePossession(
		match.HomeTeam.TeamAbility,
		match.AwayTeam.TeamAbility,
	)
	zap.L().Info("Home Possession", zap.Any("value", homePoss))
	return &models.MatchResult{
		HomeTeam: match.HomeTeam,
		AwayTeam: match.AwayTeam,
		Referee:  match.Referee,
		RngSeed:  rng.GetSeed(),
		Possession: [2]int{
			int(homePoss),
			100 - int(homePoss),
		},
	}, nil
}
