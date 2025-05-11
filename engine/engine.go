package engine

import (
	"fmt"
	"math"

	"github.com/esmshub/esms-go/engine/commentary"
	"github.com/esmshub/esms-go/engine/events"
	"github.com/esmshub/esms-go/engine/flows"
	"github.com/esmshub/esms-go/engine/models"
	"github.com/esmshub/esms-go/engine/pkg/rng"
	"github.com/esmshub/esms-go/engine/pkg/utils"
	"github.com/esmshub/esms-go/engine/validators"
	"go.uber.org/zap"
)

// Function to calculate possession
func calcHomePossession(homeTeam, awayTeam *models.MatchTeam) float64 {
	// Normalize the stats for each team
	homeAbility := homeTeam.GetAbility()
	awayAbility := awayTeam.GetAbility()
	normalizedPassingHome := float64(homeAbility.Passing) / MaxTeamSkillValue
	normalizedTacklingHome := float64(homeAbility.Tackling) / MaxTeamSkillValue
	normalizedShootingHome := float64(homeAbility.Shooting) / MaxTeamSkillValue

	normalizedPassingAway := float64(awayAbility.Tackling) / MaxTeamSkillValue
	normalizedTacklingAway := float64(awayAbility.Passing) / MaxTeamSkillValue
	normalizedShootingAway := float64(awayAbility.Shooting) / MaxTeamSkillValue

	// Calculate the initial possession
	totalStrength := normalizedPassingHome + normalizedTacklingHome + normalizedShootingHome +
		normalizedPassingAway + normalizedTacklingAway + normalizedShootingAway

	poss := (normalizedPassingHome + normalizedTacklingHome + normalizedShootingHome) / totalStrength * 100

	// Add a random factor for variability
	randomFactor := (0.05 + rng.RandomF()*0.10) // A random factor between 0.05 and 0.15
	return math.Min(math.Floor(poss*(1+randomFactor)), 85.0)
}

func Run(match *models.Match, options *Options) (*models.MatchResult, error) {
	if options.RngSeed != 0 {
		// seed random number generator
		zap.L().Info("Seeding RNG", zap.Uint64("seed", options.RngSeed))
		rng.Seed(options.RngSeed)
	}
	if match.Referee == nil {
		zap.L().Warn("No referee assigned, using default")
		match.Referee = models.NewDefaultReferee()
	}

	validator := validators.NewTeamConfigValidator(options.AppConfig)
	for _, team := range match.GetTeams() {
		if err := validator.Validate(team); err != nil {
			return nil, err
		}
	}

	fmt.Println("---------- Kick off ----------")
	zap.L().Info("Initial Team ability", zap.Any("home", match.HomeTeam.GetAbility()))
	zap.L().Info("Initial Team ability", zap.Any("away", match.AwayTeam.GetAbility()))
	zap.L().Info("Before probability", zap.Any("shot_prob", match.HomeTeam.GetShotProbability()))

	absCalculator := models.NewAbilityCalculator(
		options.TacticsMatrix,
		DefAggressionLevel,
	)
	fatigueCalculator := &models.MatchFatigueCalculator{}
	probabilityCalculator := &models.ProbabilityCalculator{}
	statsUpdater := &models.MatchStatsUpdater{}
	// hook up observers
	match.Subscribe(fatigueCalculator)
	match.Subscribe(absCalculator)
	match.Subscribe(statsUpdater)
	match.Subscribe(probabilityCalculator)
	defer match.UnsubscribeAll()

	// Create event bus
	eventBus := events.NewMemoryEventBus()
	eventBus.RegisterHandler(models.MatchEventHandler{})
	eventBus.RegisterHandler(commentary.NewEventHandler(options.CommentaryProvider))
	// eventBus.RegisterHandler(events.ShotEventType, NewShotHandler(match))
	// eventBus.RegisterHandler(events.ShotOnTargetType, NewShotOnTargetHandler(match))
	// ... register other handlers ...

	// main loop
	flowIndices := utils.Range(0, len(flows.All))
	teams := match.GetTeams()
	teamIndices := []int{0, 1}
	stats := []models.MatchStats{}
	for {
		// randomize flow order
		rng.Shuffle(len(flowIndices), func(i, j int) {
			flowIndices[i], flowIndices[j] = flowIndices[j], flowIndices[i]
		})

		// randomize team selection
		rng.Shuffle(len(teamIndices), func(i, j int) {
			teamIndices[i], teamIndices[j] = teamIndices[j], teamIndices[i]
		})

		if match.GetMinute() == 0 {
			// kick off event
			eventBus.Publish(models.NewMatchEvent(events.NewEvent(events.KickOffEventName, nil), match, teams[0]))
		}

		for _, i := range teamIndices {
			zap.L().Debug("Team selected", zap.Any("team", teams[i].GetName()))
			for _, x := range flowIndices {
				team := teams[i]
				if team == nil {
					zap.L().Panic("team not found", zap.Int("index", i))
				}
				flow := flows.All[x]
				if flow == nil {
					zap.L().Panic("flow not found", zap.Int("index", x))
				}

				flowEvents := flow(team, teams[i^1])
				zap.L().Debug("Flow events", zap.Any("events", flowEvents))

				utils.EachFunc(flowEvents, func(event events.Event) {
					eventBus.Publish(models.NewMatchEvent(event, match, team))
				})
			}
		}

		match.IncrementMinute()

		curStats := models.MatchStats{
			InjuryTime: match.GetInjuryTime(),
			HomeStats:  match.HomeTeam.GetStats(),
			AwayStats:  match.AwayTeam.GetStats(),
		}
		curMinute := match.GetMinute()
		isFirstHalf := len(stats) == 0
		if isFirstHalf && curMinute == MinsPerHalf+curStats.InjuryTime {
			stats = append(stats, curStats)
			eventBus.Publish(models.NewMatchEvent(events.NewEvent(events.HalfTimeEventName, nil), match, nil))
			match.SetMinute(MinsPerHalf)
		} else if !isFirstHalf && curMinute == (MinsPerHalf*2)+curStats.InjuryTime {
			stats = append(stats, curStats)
			eventBus.Publish(models.NewMatchEvent(events.NewEvent(events.FullTimeEventName, nil), match, nil))
			break
		} else if curMinute == MinsPerHalf-2 || curMinute == (MinsPerHalf*2)-2 {
			match.CalculateInjuryTime()
			eventBus.Publish(models.NewMatchEvent(events.NewEvent(events.InjuryTimeAddedEventName, map[string]any{"injury_time": match.GetInjuryTime()}), match, nil))
		}
	}

	halfTimeStats := stats[0]
	fmt.Println("---------- Half time ----------")
	zap.L().Info("Half time stats", zap.Any("stats", halfTimeStats))
	zap.L().Info("After Team ability", zap.Any("home", match.HomeTeam.GetAbility()))
	zap.L().Info("After Team ability", zap.Any("away", match.AwayTeam.GetAbility()))
	zap.L().Info("After Shot probability", zap.Any("shot_prob", match.HomeTeam.GetShotProbability()))

	fullTimeStats := stats[1]
	fmt.Println("---------- Full time ----------")
	zap.L().Info("Full time stats", zap.Any("stats", fullTimeStats))
	zap.L().Info("After Team ability", zap.Any("home", match.HomeTeam.GetAbility()))
	zap.L().Info("After Team ability", zap.Any("away", match.AwayTeam.GetAbility()))
	zap.L().Info("After Shot probability", zap.Any("shot_prob", match.HomeTeam.GetShotProbability()))

	// calc possession
	homePoss := calcHomePossession(match.HomeTeam, match.AwayTeam)
	match.HomeTeam.SetPossession(int(homePoss))
	match.AwayTeam.SetPossession(100 - int(homePoss))
	zap.L().Info("Home Possession", zap.Any("value", homePoss))
	return &models.MatchResult{
		HomeTeam: match.HomeTeam,
		AwayTeam: match.AwayTeam,
		Referee:  match.Referee,
		RngSeed:  rng.GetSeed(),
	}, nil
}
