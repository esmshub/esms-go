package models

import (
	"math/rand"
	"slices"

	"github.com/esmshub/esms-go/engine/pkg/utils"
	"github.com/esmshub/esms-go/engine/types"
	"go.uber.org/zap"
)

type MatchBonusCalculator struct {
	bonuses map[string]any
}

func (b *MatchBonusCalculator) getBonus(key string) int {
	v, exists := b.bonuses[key]
	if !exists {
		zap.L().Warn("bonus not set", zap.String("key", key))
		return 0
	}

	if result, ok := v.(int); ok {
		return result
	} else {
		zap.L().Warn("bonus not an int", zap.String("key", key))
		return 0
	}
}

func (b *MatchBonusCalculator) Apply(result *MatchResult) {
	teams := []*MatchTeam{result.HomeTeam, result.AwayTeam}
	for i, t := range teams {
		for _, p := range t.GetLineup() {
			abs := p.GetMatchAbility()
			// shot stopping
			abs.GoalkeepingAbs = b.getBonus("ab_sav") * p.GetStats().Saves
			abs.GoalkeepingAbs += b.getBonus("ab_concede") * p.GetStats().Conceded
			// tackling
			abs.TacklingAbs = b.getBonus("ab_ktk") * p.GetStats().KeyTackles
			abs.TacklingAbs += b.getBonus("ab_og") * len(p.GetStats().OwnGoals)
			// passing
			abs.PassingAbs = b.getBonus("ab_kps") * p.GetStats().KeyPasses
			abs.PassingAbs += b.getBonus("ab_assist") * p.GetStats().Assists
			abs.PassingAbs += b.getBonus("ab_og") * len(p.GetStats().OwnGoals)
			// shooting
			abs.ShootingAbs = b.getBonus("ab_goal") * len(p.GetStats().Goals)
			abs.ShootingAbs += b.getBonus("ab_sht_on") * (p.GetStats().ShotsOffTarget)
			abs.ShootingAbs += b.getBonus("ab_sht_off") * p.GetStats().ShotsOffTarget

			if p.GetStats().IsCautioned {
				cautionedBonus := b.getBonus("ab_yellow")
				if p.GetPosition() == types.PositionGK {
					abs.GoalkeepingAbs += cautionedBonus
				} else {
					abs.TacklingAbs += cautionedBonus
					abs.PassingAbs += cautionedBonus
					abs.ShootingAbs += cautionedBonus
				}
			}
			if p.GetStats().IsSentOff {
				sentOffBonus := b.getBonus("ab_red")
				if p.GetPosition() == types.PositionGK {
					abs.GoalkeepingAbs += sentOffBonus
				} else {
					abs.TacklingAbs += sentOffBonus
					abs.PassingAbs += sentOffBonus
					abs.ShootingAbs += sentOffBonus
				}
			}

			p.SetMatchAbility(&abs)
		}

		if len(t.GetStats().Goals) > len(teams[i^1].GetStats().Goals) {
			b.applyVictoryBonus(t)
		} else if len(t.GetStats().Goals) < len(teams[i^1].GetStats().Goals) {
			b.applyDefeatBonus(t)
		}

		if len(teams[i^1].GetStats().Goals) == 0 {
			b.applyCleanSheetBonus(t)
		}
	}
}

func (b *MatchBonusCalculator) applyVictoryBonus(team *MatchTeam) {
	zap.L().Debug("applying victory bonus", zap.String("team", team.GetName()))
	victoryBonus := b.getBonus("ab_victory_random")
	if victoryBonus == 0 {
		return
	}

	players := b.getRandomPlayers(team, 2, func(p *MatchPlayer) bool {
		return p.GetStats().MinutesPlayed > 0
	})
	for _, p := range players {
		abs := p.GetMatchAbility()
		if p.GetPosition() == types.PositionGK {
			abs.GoalkeepingAbs += victoryBonus
		} else {
			abs.TacklingAbs += victoryBonus
			abs.PassingAbs += victoryBonus
			abs.ShootingAbs += victoryBonus
		}
		p.SetMatchAbility(&abs)
	}
}

func (b *MatchBonusCalculator) applyDefeatBonus(team *MatchTeam) {
	zap.L().Debug("applying defeat bonus", zap.String("team", team.GetName()))
	defeatBonus := b.getBonus("ab_defeat_random")
	if defeatBonus == 0 {
		return
	}

	players := b.getRandomPlayers(team, 2, func(p *MatchPlayer) bool {
		return p.GetStats().MinutesPlayed > 0
	})
	for _, p := range players {
		abs := p.GetMatchAbility()
		if p.GetPosition() == types.PositionGK {
			abs.GoalkeepingAbs += defeatBonus
		} else {
			abs.TacklingAbs += defeatBonus
			abs.PassingAbs += defeatBonus
			abs.ShootingAbs += defeatBonus
		}
		p.SetMatchAbility(&abs)
	}
}

func (b *MatchBonusCalculator) applyCleanSheetBonus(team *MatchTeam) {
	zap.L().Debug("applying clean sheet bonus", zap.String("team", team.GetName()))
	cleanSheetBonus := b.getBonus("ab_clean_sheet")
	if cleanSheetBonus == 0 {
		return
	}

	keeper := utils.FindFunc(team.GetLineup(), func(p *MatchPlayer) bool {
		return p.GetPosition() == types.PositionGK && p.GetStats().MinutesPlayed > 45
	})
	if keeper == nil {
		zap.L().Warn("no keeper with 45+ minutes played")
		return
	}
	// apply bonus to keeper
	keeperAbs := keeper.GetMatchAbility()
	keeperAbs.GoalkeepingAbs += cleanSheetBonus
	keeper.SetMatchAbility(&keeperAbs)

	// apply bonus to defender
	defender := b.getRandomPlayers(team, 1, func(p *MatchPlayer) bool {
		return p.GetPosition() == types.PositionDF && p.GetStats().MinutesPlayed > 0
	})[0]
	defenderAbs := defender.GetMatchAbility()
	defenderAbs.TacklingAbs += cleanSheetBonus
	defender.SetMatchAbility(&defenderAbs)
}

func (b *MatchBonusCalculator) getRandomPlayers(team *MatchTeam, count int, conditionFunc func(p *MatchPlayer) bool) []*MatchPlayer {
	players := team.GetLineup()
	indexes := []int{}
	for len(indexes) < count {
		i := rand.Intn(len(players))
		if conditionFunc(players[i]) && !slices.Contains(indexes, i) {
			indexes = append(indexes, i)
		}
	}

	return utils.Map(indexes, func(i int) *MatchPlayer {
		return players[i]
	})
}

func NewMatchBonusCalculator(bonuses map[string]any) *MatchBonusCalculator {
	return &MatchBonusCalculator{
		bonuses: bonuses,
	}
}
