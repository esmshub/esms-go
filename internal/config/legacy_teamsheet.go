package config

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/esmshub/esms-go/engine"
	"github.com/esmshub/esms-go/engine/models"
	"go.uber.org/zap"
)

func readTactic(config *models.TeamConfig, text string) error {
	zap.L().Debug("reading tactic")
	lines := strings.Split(text, "\n")
	if len(lines) > 2 {
		return errors.New("unexpected number of rows in first section")
	}

	if len(lines) > 1 {
		config.Name = lines[0]
		config.Tactic = lines[1]
	} else {
		config.Tactic = lines[0]
	}

	if _, exists := engine.TacticNames[config.Tactic]; !exists {
		return fmt.Errorf("tactic '%s' is not a recognised value", config.Tactic)
	}

	return nil
}

func parsePlayer(text string, playerIndex int, isSub bool, findPlayer func(string) *models.Player) (*models.Player, error) {
	parts := strings.Fields(text)
	if len(parts) != 2 {
		return nil, fmt.Errorf("player data format '%s' is not valid", text)
	}
	pos := strings.TrimSpace(parts[0])
	name := strings.TrimSpace(parts[1])

	zap.L().Debug("validate player position", zap.String("pos", pos))
	if !IsValidPosition(pos) {
		return nil, fmt.Errorf("player %s has an illegal position (%s)", name, pos)
	} else if playerIndex == 1 && pos != POSITION_GK {
		return nil, fmt.Errorf("player 1 must be a GK")
	} else if pos == POSITION_GK && playerIndex != 1 && !isSub {
		return nil, fmt.Errorf("player %s cannot be a GK", name)
	}

	zap.L().Debug("searching player in roster", zap.String("name", name))
	if player := findPlayer(name); player != nil {
		zap.L().Debug("player found")
		if player.GetIsInjured() {
			return nil, fmt.Errorf("player %s is injured", name)
		} else if player.GetIsSuspended() {
			return nil, fmt.Errorf("player %s is suspended", name)
		}

		player.Position = pos
		return player, nil
	} else {
		return nil, fmt.Errorf("player %s does not exist in the roster", name)
	}
}

func readLineup(config *models.TeamConfig, text string, findPlayer func(string) *models.Player) error {
	zap.L().Debug("reading lineup")

	config.Lineup = []*models.Player{}
	lines := strings.Split(text, "\n")
	playerCount := len(lines)
	if playerCount != 11 {
		return fmt.Errorf("expected 11 first team players, found %d", len(lines))
	}

	for idx, record := range lines {
		player, err := parsePlayer(strings.TrimSpace(record), idx+1, false, findPlayer)
		if err != nil {
			return err
		}

		config.Lineup = append(config.Lineup, player)
	}
	return nil
}

func readSubstitutions(config *models.TeamConfig, text string, minSubs int, maxSubs int, findPlayer func(string) *models.Player) error {
	zap.L().Debug("reading subs")

	config.Subs = []*models.Player{}
	lines := strings.Split(text, "\n")
	subCount := len(lines)
	if subCount < minSubs {
		return fmt.Errorf("expected at least %d subs, found %d", minSubs, subCount)
	} else if subCount > maxSubs {
		return fmt.Errorf("expected at most %d subs, found %d", maxSubs, subCount)
	}

	for idx, record := range lines {
		player, err := parsePlayer(strings.TrimSpace(record), idx+1, true, findPlayer)
		if err != nil {
			return err
		}

		config.Subs = append(config.Subs, player)
	}
	return nil
}

func readPenaltyTaker(config *models.TeamConfig, text string) error {
	zap.L().Debug("reading pk taker")
	lines := strings.Fields(text)
	if len(lines) != 2 {
		return fmt.Errorf("invalid PK taker")
	}

	pk := strings.TrimSpace(lines[1])
	nameComparer := func(p *models.Player) bool {
		return strings.EqualFold(p.Name, pk)
	}
	i := slices.IndexFunc(config.Lineup, nameComparer)
	if i == -1 {
		return fmt.Errorf("PK taker must be in the starting lineup")
	}

	config.PlayerRoles["PK"] = config.Lineup[i]
	return nil
}

func validateFormation(team *models.TeamConfig) error {
	positionCount := make(map[string]int)
	for _, p := range team.Lineup {
		positionCount[p.Position]++
	}
	defMids := positionCount[POSITION_DM]
	mids := positionCount[POSITION_MF]
	attMids := positionCount[POSITION_AM]

	if positionCount[POSITION_DF] < LeagueConfig.GetInt("min_df") {
		return fmt.Errorf("cannot field less than %d %ss", LeagueConfig.GetInt("min_df"), POSITION_DF)
	} else if positionCount[POSITION_DF] > LeagueConfig.GetInt("max_df") {
		return fmt.Errorf("cannot field more than %d %ss", LeagueConfig.GetInt("max_df"), POSITION_DF)
	} else if positionCount[POSITION_DM] > LeagueConfig.GetInt("max_dm") {
		return fmt.Errorf("cannot field more than %d %ss", LeagueConfig.GetInt("max_dm"), POSITION_DM)
	} else if positionCount[POSITION_MF] > LeagueConfig.GetInt("max_mf") {
		return fmt.Errorf("cannot field more than %d %ss", LeagueConfig.GetInt("max_mf"), POSITION_MF)
	} else if positionCount[POSITION_AM] > LeagueConfig.GetInt("max_am") {
		return fmt.Errorf("cannot field more than %d %ss", LeagueConfig.GetInt("max_am"), POSITION_AM)
	} else if (mids + defMids + attMids) < LeagueConfig.GetInt("min_mf") {
		return fmt.Errorf("cannot field less than %d midfielders", LeagueConfig.GetInt("min_mf"))
	} else if (mids + defMids + attMids) > LeagueConfig.GetInt("max_mf") {
		return fmt.Errorf("cannot field more than %d midfielders", LeagueConfig.GetInt("max_mf"))
	} else if positionCount[POSITION_FW] < LeagueConfig.GetInt("min_fw") {
		return fmt.Errorf("cannot field less than %d %ss", LeagueConfig.GetInt("min_fw"), POSITION_FW)
	} else if positionCount[POSITION_FW] > LeagueConfig.GetInt("max_fw") {
		return fmt.Errorf("cannot field more than %d %ss", LeagueConfig.GetInt("max_fw"), POSITION_FW)
	}

	formationStr := fmt.Sprintf("%d-%d-%d-%d-%d", positionCount[POSITION_DF], defMids, mids, attMids, positionCount[POSITION_FW])
	team.Formation = strings.ReplaceAll(formationStr, "-0", "")

	return nil
}

func readConditionals(config *models.TeamConfig, text string) error {
	config.Conditionals = []*models.Conditional{}
	lines := strings.Split(string(text), "\n")
	for _, l := range lines {
		cond := strings.ToUpper(strings.TrimSpace(l))
		var parse ConditionalParser = nil
		if strings.HasPrefix(cond, AggressionAction) {
			parse = conditionalParsers[AggressionAction]
		} else if strings.HasPrefix(cond, ChangeAggressionAction) {
			parse = conditionalParsers[ChangeAggressionAction]
		} else if strings.HasPrefix(cond, ChangePositionAction) {
			parse = conditionalParsers[ChangePositionAction]
		} else if strings.HasPrefix(cond, SubstitutionAction) {
			parse = conditionalParsers[SubstitutionAction]
		} else if strings.HasPrefix(cond, TacticAction) {
			parse = conditionalParsers[TacticAction]
		} else {
			return fmt.Errorf("unknown conditional: %s", cond)
		}

		if parse == nil {
			panic(fmt.Errorf("no parser for %s conditional found", cond))
		}

		conditional, err := parse(cond)
		if err != nil {
			return err
		}

		config.Conditionals = append(config.Conditionals, conditional)
	}
	return nil
}

func LoadLegacyTeamsheet(teamsheetFile string, findPlayer func(string) *models.Player) (*models.TeamConfig, error) {
	zap.L().Info("reading teamsheet file", zap.String("file", teamsheetFile))
	data, err := os.ReadFile(teamsheetFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read teamsheet: %w", err)
	}

	normalisedFileContents := strings.ReplaceAll(string(data), "\r\n", "\n")
	blocks := strings.Split(normalisedFileContents, "\n\n")
	sections := []string{}
	for _, block := range blocks {
		if len(strings.TrimSpace(block)) > 0 {
			// fmt.Println("-----------")
			// fmt.Println(block)
			sections = append(sections, strings.TrimSpace(block))
		}
	}

	if len(sections) < 4 {
		return nil, errors.New("invalid teamsheet: missing sections")
	}

	config := &models.TeamConfig{
		PlayerRoles: make(map[string]*models.Player),
	}
	if err := readTactic(config, sections[0]); err != nil {
		return nil, err
	}

	if err := readLineup(config, sections[1], findPlayer); err != nil {
		return nil, err
	}
	if err := readSubstitutions(config, sections[2], LeagueConfig.GetInt("min_subs"), LeagueConfig.GetInt("bench_size"), findPlayer); err != nil {
		return nil, err
	}
	if err := readPenaltyTaker(config, sections[3]); err != nil {
		return nil, err
	}
	if err := validateFormation(config); err != nil {
		return nil, err
	}

	if len(sections) > 4 {
		if err := readConditionals(config, sections[4]); err != nil {
			return nil, err
		}
	}

	return config, err
}
