package formatters

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/esmshub/esms-go/engine/models"
	"github.com/esmshub/esms-go/engine/pkg/utils"
	"github.com/esmshub/esms-go/engine/types"
)

func getMatchData(event *models.MatchEvent) *map[string]any {
	if data, ok := event.GetData().(map[string]any); ok {
		return &data
	}
	return &map[string]any{}
}

func FormatMatchEvent(format string, event *models.MatchEvent, args ...any) string {
	result := format
	_args := []any{}
	if strings.Contains(result, "Min. %s") {
		_args = append(_args, strconv.Itoa(event.GetMatch().GetMinute()))
	} else if strings.Contains(result, "Min. %-3d") || strings.Contains(result, "Min. %d") {
		_args = append(_args, event.GetMatch().GetMinute())
	}

	if event.GetActiveTeam() != nil {
		if strings.Contains(result, "(%s)") {
			_args = append(_args, event.GetActiveTeam().GetShortName())
		}
		result = strings.ReplaceAll(result, "<m1>", event.GetActiveTeam().GetManagerName())
		result = strings.ReplaceAll(result, "<t1>", event.GetActiveTeam().GetName())
		result = strings.ReplaceAll(result, "<g1>", event.GetActiveTeam().GetFirstActiveByPosition(types.PositionGK).GetName())
		result = strings.ReplaceAll(result, "<m2>", event.GetOppositionTeam().GetManagerName())
		result = strings.ReplaceAll(result, "<t2>", event.GetOppositionTeam().GetName())
		result = strings.ReplaceAll(result, "<g2>", event.GetOppositionTeam().GetFirstActiveByPosition(types.PositionGK).GetName())
	}

	result = strings.ReplaceAll(result, "<ref>", event.GetMatch().Referee.Name)
	result = strings.ReplaceAll(result, "<venue>", event.GetMatch().HomeTeam.GetStadiumName())

	_args = append(_args, args...)
	// zap.L().Warn(event.GetName(), zap.Any("result", result), zap.Any("_args", _args))
	return fmt.Sprintf(result, _args...)
}

func FormatKickOffEvent(format string, event *models.MatchEvent) string {
	return FormatMatchEvent("Min. %d :(REF) There's the whistle, %s to kick off and get the game underway", event, event.GetActiveTeam().GetName())
}

func FormatInjuryTimeEvent(format string, event *models.MatchEvent) string {
	data := *getMatchData(event)
	mins := data["injury_time"].(int)
	if mins == 0 {
		return FormatMatchEvent("\nMin. %d :(REF) The 4th official keeps the board down, doesn't look like there's going to be any added time.", event)
	}
	return FormatMatchEvent(format, event, mins)
}

func FormatHalfTimeEvent(format string, event *models.MatchEvent) string {
	return FormatMatchEvent("\nMin. %d :(REF) The whistle signals the end of the first half.\n%s\n", event, format)
}

func FormatFullTimeEvent(format string, event *models.MatchEvent) string {
	return FormatMatchEvent("\nMin. %d :(REF) The referee brings the game to an end with the final whistle.\n%s\n", event, format)
}

func FormatScoreEvent(format string, event *models.MatchEvent) string {
	match := event.GetMatch()
	homeGoals := len(match.HomeTeam.GetStats().Goals)
	awayGoals := len(match.AwayTeam.GetStats().Goals)
	return FormatMatchEvent(format, event, match.HomeTeam.GetName(), homeGoals, awayGoals, match.AwayTeam.GetName())
}

func FormatAssistedShotEvent(format string, event *models.MatchEvent) string {
	data := *getMatchData(event)
	return FormatMatchEvent(
		format,
		event,
		data["assister"].(*models.MatchPlayer).GetName(),
		data["attacker"].(*models.MatchPlayer).GetName(),
	)
}

func FormatShotEvent(format string, event *models.MatchEvent) string {
	data := *getMatchData(event)
	attacker := utils.MustGetKey[*models.MatchPlayer](data, "attacker")
	args := []any{attacker.GetName()}
	if oneOnOne, ok := data["one_on_one"].(bool); ok && oneOnOne {
		keeper := utils.MustGetKey[*models.MatchPlayer](data, "opp_keeper")
		args = append(args, keeper.GetName())
	}
	return FormatMatchEvent(format, event, args...)
}

func FormatShotTackledEvent(format string, event *models.MatchEvent) string {
	data := *getMatchData(event)
	return FormatMatchEvent(format, event, data["tackler"].(*models.MatchPlayer).GetName())
}

func FormatShotSavedEvent(format string, event *models.MatchEvent) string {
	data := *getMatchData(event)
	keeper := data["opp_keeper"].(*models.MatchPlayer)
	return FormatMatchEvent(format, event, keeper.GetName())
}

func FormatShotOffTargetEvent(format string, event *models.MatchEvent) string {
	return FormatMatchEvent(format, event)
}

func FormatChanceEvent(format string, event *models.MatchEvent) string {
	data := *getMatchData(event)
	attacker := data["attacker"].(*models.MatchPlayer)
	args := []any{attacker.GetName()}
	if def, ok := data["got_past_defender"].(*models.MatchPlayer); ok {
		args = append(args, def.GetName())
	}
	return FormatMatchEvent(format, event, args...)
}

func FormatAssistedChanceEvent(format string, event *models.MatchEvent) string {
	data := *getMatchData(event)
	assister := data["assister"].(*models.MatchPlayer)
	args := []any{assister.GetName()}
	if def, ok := data["got_past_defender"].(*models.MatchPlayer); ok {
		args = append(args, def.GetName())
	}
	attacker := data["attacker"].(*models.MatchPlayer)
	args = append(args, attacker.GetName())
	return FormatMatchEvent(format, event, args...)
}

func FormatOwnGoalScoredEvent(format string, event *models.MatchEvent) string {
	data := *getMatchData(event)
	return FormatMatchEvent(format, event, data["scorer"].(*models.MatchPlayer).GetName())
}

func FormatCornerEvent(format string, event *models.MatchEvent) string {
	data := *getMatchData(event)
	return FormatMatchEvent(format, event, data["corner_taker"].(*models.MatchPlayer).GetName())
}

func FormatCornerCaughtEvent(format string, event *models.MatchEvent) string {
	keeper := event.GetOppositionTeam().GetActiveByPosition(types.PositionGK)[0]
	return FormatMatchEvent(format, event, keeper.GetName())
}

func FormatCornerClearedEvent(format string, event *models.MatchEvent) string {
	data := *getMatchData(event)
	return FormatMatchEvent(format, event, data["tackler"].(*models.MatchPlayer).GetName())
}

func FormatCornerShotEvent(format string, event *models.MatchEvent) string {
	data := *getMatchData(event)
	return FormatMatchEvent(format, event, data["attacker"].(*models.MatchPlayer).GetName())
}

func FormatYellowCardEvent(format string, event *models.MatchEvent) string {
	return FormatMatchEvent(format, event)
}

func FormatRedCardEvent(format string, event *models.MatchEvent) string {
	return FormatMatchEvent(format, event)
}

func FormatSubstitutionEvent(format string, event *models.MatchEvent) string {
	return FormatMatchEvent(format, event)
}

func FormatTackleEvent(format string, event *models.MatchEvent) string {
	return FormatMatchEvent(format, event)
}
