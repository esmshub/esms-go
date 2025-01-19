package store

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/esmshub/esms-go/engine/models"
	"github.com/esmshub/esms-go/pkg/utils"
	"go.uber.org/zap"
)

type MatchResultFileStoreOptions struct {
	OutputFile string
	HeaderText string
	FooterText string
}

type MatchResultFileStore struct{}

func (t *MatchResultFileStore) saveAsText(result *models.MatchResult, options MatchResultFileStoreOptions) error {
	zap.L().Info("saving match result as text", zap.String("file", options.OutputFile))
	// Open the file in write mode. If the file doesn't exist, it will be created.
	file, err := os.Create(options.OutputFile)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close() // Make sure to close the file when done

	// Create a new buffered writer for the file
	writer := bufio.NewWriter(file)

	lines := []string{}
	// print match info
	now := time.Now()
	if options.HeaderText != "" {

	} else {

	}
	lines = append(lines, fmt.Sprintf("%s, %s vs. %s (%s)\n", options.HeaderText, result.HomeTeam.Name, result.AwayTeam.Name, now.Format("Mon Jan 02")))
	lines = append(lines, fmt.Sprintf("%37s  |  %s", result.HomeTeam.Name, result.AwayTeam.Name))
	lines = append(lines, fmt.Sprintf("%37s  |  %s", result.HomeTeam.Formation, result.AwayTeam.Formation))
	lines = append(lines, fmt.Sprintf("%37s  |  %s", result.HomeTeam.ManagerName, result.AwayTeam.ManagerName))
	lines = append(lines, fmt.Sprintf("%37s  |", ""))
	// print lineup
	for i := 0; i < len(result.HomeTeam.Lineup); i++ {
		hp := result.HomeTeam.Lineup[i]
		ap := result.AwayTeam.Lineup[i]
		lines = append(lines, fmt.Sprintf("%34s %s  |  %s %s", hp.Name, hp.Position, ap.Position, ap.Name))
	}
	lines = append(lines, fmt.Sprintf("%37s  |", ""))
	// print subs
	for i := 0; i < int(math.Max(float64(len(result.HomeTeam.Subs)), float64(len(result.AwayTeam.Subs)))); i++ {
		hs, as := "", ""
		if i < len(result.HomeTeam.Subs) {
			hs = fmt.Sprintf("%s SUB", result.HomeTeam.Subs[i].Name)
		}
		if i < len(result.AwayTeam.Subs) {
			as = fmt.Sprintf("SUB %s", result.AwayTeam.Subs[i].Name)
		}
		lines = append(lines, fmt.Sprintf("%37s  |  %s", hs, as))
	}
	// TODO: print ref
	if result.Referee != nil {
		lines = append(lines, fmt.Sprintf("\nReferee: %s (%s)", result.Referee.Name, result.Referee.Nat))
	}
	// TODO: print commentary
	lines = append(lines, "\n")
	// print match details
	lines = append(lines, "Match Details")
	lines = append(lines, strings.Repeat("-", 91))
	lines = append(lines, fmt.Sprintf("%-22s: %s", "Venue", result.HomeTeam.StadiumName))
	lines = append(lines, fmt.Sprintf("%-22s: %d", "Attendance", 0))
	lines = append(lines, strings.Repeat("-", 91))
	// TODO: attendance
	// print match info (home)
	lines = append(lines, result.HomeTeam.Name+" Match Info")
	lines = append(lines, strings.Repeat("-", 91))
	lines = append(lines, fmt.Sprintf("%-22s: %s %s", "Best player", result.HomeTeam.Lineup[4].Name, "(Man of the Match)"))
	lines = append(lines, fmt.Sprintf("%-22s: %s", "Scorers", "T_Kierney (6)"))
	lines = append(lines, fmt.Sprintf("%-22s: %s", "Sent Off", "N/A"))
	lines = append(lines, fmt.Sprintf("%-22s: %s", "Booked", "N/A"))
	lines = append(lines, fmt.Sprintf("%-22s: %s", "Injured", "N/A"))
	lines = append(lines, strings.Repeat("-", 91))
	// print match info (away)
	lines = append(lines, result.AwayTeam.Name+" Match Info")
	lines = append(lines, strings.Repeat("-", 91))
	lines = append(lines, fmt.Sprintf("%-22s: %s %s", "Best player", result.AwayTeam.Lineup[4].Name, "(Man of the Match)"))
	lines = append(lines, fmt.Sprintf("%-22s: %s", "Scorers", "K_Tierney (12)"))
	lines = append(lines, fmt.Sprintf("%-22s: %s", "Sent Off", "N/A"))
	lines = append(lines, fmt.Sprintf("%-22s: %s", "Booked", "N/A"))
	lines = append(lines, fmt.Sprintf("%-22s: %s", "Injured", "N/A"))
	lines = append(lines, strings.Repeat("-", 91))
	// print match stats
	lines = append(lines, "\nMatch Statistics")
	lines = append(lines, strings.Repeat("-", 91))
	lines = append(lines, fmt.Sprintf("%40s   |    %s", result.HomeTeam.Name, result.AwayTeam.Name))
	lines = append(lines, fmt.Sprintf("%-22s:         %3d        |       %d", "Attempts on Goal", 1, 1))
	lines = append(lines, fmt.Sprintf("%-22s:         %3d        |       %d", "Blocked Shots", 1, 1))
	lines = append(lines, fmt.Sprintf("%-22s:         %3d        |       %d", "Off target", 1, 1))
	lines = append(lines, fmt.Sprintf("%-22s:         %3d        |       %d", "On target", 1, 1))
	lines = append(lines, fmt.Sprintf("%-22s:         %3d        |       %d", "Possession", result.Possession[0], result.Possession[1]))
	lines = append(lines, fmt.Sprintf("%-22s:         %3d        |       %d", "Corners", 1, 1))
	lines = append(lines, fmt.Sprintf("%-22s:         %3d        |       %d", "Red Cards", 1, 1))
	lines = append(lines, fmt.Sprintf("%-22s:         %3d        |       %d", "Yellow Cards", 1, 1))
	lines = append(lines, fmt.Sprintf("%-22s:         %3d        |       %d", "Subs Used", 1, 1))
	lines = append(lines, fmt.Sprintf("%-22s:         %3d        |       %d", "Fouls", 1, 1))
	lines = append(lines, fmt.Sprintf("%-22s:         %3d        |       %d", "Key Tackles", 1, 1))
	lines = append(lines, fmt.Sprintf("%-22s:         %3d        |       %d", "Key Passes", 1, 1))
	lines = append(lines, strings.Repeat("-", 91))
	lines = append(lines, fmt.Sprintf("%-22s: %s %d - %d %s", "Final Score", result.HomeTeam.Name, 0, 0, result.AwayTeam.Name))
	lines = append(lines, strings.Repeat("-", 91))
	// print player stats (home)
	lines = append(lines, "\nPlayer statistics - "+result.HomeTeam.Name)
	lines = append(lines, "Name          Pos  St  Tk  Ps  Sh  Ag | Min Sav Ktk Kps Ass Sht Gls Yel Red KAb TAb PAb SAb")
	lines = append(lines, strings.Repeat("-", 91))
	homeSquad := append(result.HomeTeam.Lineup, result.HomeTeam.Subs...)
	for _, p := range homeSquad {
		lines = append(lines, fmt.Sprintf("%-13s %3s %3d %3d %3d %3d %3d | %3d %3d %3d %3d %3d %3d %3d %3d %3d %3d %3d %3d %3d",
			p.Name,
			p.Position,
			p.BaseAbility.Goalkeeping,
			p.BaseAbility.Tackling,
			p.BaseAbility.Passing,
			p.BaseAbility.Shooting,
			p.BaseAbility.Aggression,
			p.Stats.MinutesPlayed,
			p.Stats.Saves,
			p.Stats.KeyTackles,
			p.Stats.KeyPasses,
			p.Stats.Assists,
			p.Stats.Shots,
			p.Stats.Goals,
			utils.BoolToInt(p.Stats.IsCautioned),
			utils.BoolToInt(p.Stats.IsSentOff),
			p.Ability.GoalkeepingAbs,
			p.Ability.TacklingAbs,
			p.Ability.PassingAbs,
			p.Ability.ShootingAbs,
		))
	}
	lines = append(lines, strings.Repeat("-", 91))
	// print player stats (away)
	lines = append(lines, "\nPlayer statistics - "+result.AwayTeam.Name)
	lines = append(lines, "Name          Pos  St  Tk  Ps  Sh  Ag | Min Sav Ktk Kps Ass Sht Gls Yel Red KAb TAb PAb SAb")
	lines = append(lines, strings.Repeat("-", 91))
	awaySquad := append(result.AwayTeam.Lineup, result.AwayTeam.Subs...)
	for _, p := range awaySquad {
		lines = append(lines, fmt.Sprintf("%-13s %3s %3d %3d %3d %3d %3d | %3d %3d %3d %3d %3d %3d %3d %3d %3d %3d %3d %3d %3d",
			p.Name,
			p.Position,
			p.BaseAbility.Goalkeeping,
			p.BaseAbility.Tackling,
			p.BaseAbility.Passing,
			p.BaseAbility.Shooting,
			p.BaseAbility.Aggression,
			p.Stats.MinutesPlayed,
			p.Stats.Saves,
			p.Stats.KeyTackles,
			p.Stats.KeyPasses,
			p.Stats.Assists,
			p.Stats.Shots,
			p.Stats.Goals,
			utils.BoolToInt(p.Stats.IsCautioned),
			utils.BoolToInt(p.Stats.IsSentOff),
			p.Ability.GoalkeepingAbs,
			p.Ability.TacklingAbs,
			p.Ability.PassingAbs,
			p.Ability.ShootingAbs,
		))
	}
	// print footer
	if options.FooterText != "" {
		lines = append(lines, options.FooterText)
	}
	// Write each line to the file
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n") // Write the line with a newline
		if err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	}

	// Make sure to flush the buffer to the file
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("error flushing to file: %w", err)
	}

	return nil
}

func (t *MatchResultFileStore) saveAsYaml(result *models.MatchResult, options MatchResultFileStoreOptions) error {
	return errors.ErrUnsupported
}

func (t *MatchResultFileStore) saveAsJson(result *models.MatchResult, options MatchResultFileStoreOptions) error {
	return errors.ErrUnsupported
}

func (t *MatchResultFileStore) Save(result *models.MatchResult, options MatchResultFileStoreOptions) error {
	ext := filepath.Ext(options.OutputFile)
	switch ext {
	case ".txt":
		return t.saveAsText(result, options)
	case ".yaml":
		return t.saveAsYaml(result, options)
	case ".json":
		return t.saveAsJson(result, options)
	default:
		return fmt.Errorf("unsupported file format: %s", ext)
	}
}
