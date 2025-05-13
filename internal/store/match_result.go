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
	"github.com/esmshub/esms-go/internal/formatters"
	"github.com/esmshub/esms-go/pkg/utils"
	"go.uber.org/zap"
)

type MatchResultFileStoreOptions struct {
	OutputFile string
	HeaderText string
	FooterText string
}

type MatchResultFileStore struct{}

func (mr *MatchResultFileStore) saveAsText(result *models.MatchResult, commentary []string, options MatchResultFileStoreOptions) error {
	zap.L().Info("saving match result as text", zap.String("file", options.OutputFile))
	// Open the file in write mode. If the file doesn't exist, it will be created.
	file, err := os.Create(options.OutputFile)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close() // Make sure to close the file when done

	// Create a new buffered writer for the file
	writer := bufio.NewWriter(file)

	teams := []*models.MatchTeam{result.HomeTeam, result.AwayTeam}
	lines := []string{}
	// print match info
	now := time.Now()
	if options.HeaderText != "" {

	} else {

	}
	lines = append(lines, fmt.Sprintf("%s, %s vs. %s (%s)\n", options.HeaderText, result.HomeTeam.GetName(), result.AwayTeam.GetName(), now.Format("Mon Jan 02")))
	lines = append(lines, fmt.Sprintf("%37s  |  %s", result.HomeTeam.GetName(), result.AwayTeam.GetName()))
	lines = append(lines, fmt.Sprintf("%37s  |  %s", result.HomeTeam.GetFormation(), result.AwayTeam.GetFormation()))
	lines = append(lines, fmt.Sprintf("%37s  |  %s", result.HomeTeam.GetManagerName(), result.AwayTeam.GetManagerName()))
	lines = append(lines, fmt.Sprintf("%37s  |", ""))
	// print lineup
	homeLineup := result.HomeTeam.GetStarters()
	awayLineup := result.AwayTeam.GetStarters()
	for i := 0; i < len(homeLineup); i++ {
		hp := homeLineup[i]
		ap := awayLineup[i]
		lines = append(lines, fmt.Sprintf("%34s %s  |  %s %s", hp.GetName(), hp.GetPosition(), ap.GetPosition(), ap.GetName()))
	}
	lines = append(lines, fmt.Sprintf("%37s  |", ""))
	// print subs
	homeSubs := result.HomeTeam.GetSubs()
	awaySubs := result.AwayTeam.GetSubs()
	for i := 0; i < int(math.Max(float64(len(homeSubs)), float64(len(homeSubs)))); i++ {
		hs, as := "", ""
		if i < len(homeSubs) {
			hs = fmt.Sprintf("%s SUB", homeSubs[i].GetName())
		}
		if i < len(awaySubs) {
			as = fmt.Sprintf("SUB %s", awaySubs[i].GetName())
		}
		lines = append(lines, fmt.Sprintf("%37s  |  %s", hs, as))
	}
	if result.Referee != nil {
		lines = append(lines, fmt.Sprintf("\nReferee: %s (%s)", result.Referee.Name, result.Referee.Nat))
	}
	// print commentary
	if len(commentary) > 0 {
		lines = append(lines, "\n*************  MATCH COMMENTARY  ****************")
		lines = append(lines, strings.Join(commentary, ""))
	}
	// print match details
	lines = append(lines, "\nMatch Details")
	lines = append(lines, strings.Repeat("-", 91))
	lines = append(lines, fmt.Sprintf("%-22s: %s", "Venue", result.HomeTeam.GetStadiumName()))
	lines = append(lines, fmt.Sprintf("%-22s: %d", "Attendance", 0))
	lines = append(lines, strings.Repeat("-", 91))
	// TODO: attendance
	// print match info
	for _, t := range teams {
		lines = append(lines, t.GetName()+" Match Info")
		lines = append(lines, strings.Repeat("-", 91))
		// TODO: MoM
		lines = append(lines, fmt.Sprintf("%-22s: %s", "Best player", "N/A"))
		lines = append(lines, fmt.Sprintf("%-22s: %s", "Scorers", formatters.FormatScorers(t.GetStats().Goals, formatters.FormatScorersOptions{
			RowDelimiter:         ", ",
			NoResultsPlaceholder: "N/A",
		})))
		lines = append(lines, fmt.Sprintf("%-22s: %s", "Yellow Cards", "N/A"))
		lines = append(lines, fmt.Sprintf("%-22s: %s", "Red Cards", "N/A"))
		lines = append(lines, fmt.Sprintf("%-22s: %s", "Sent Off", "N/A"))
		lines = append(lines, fmt.Sprintf("%-22s: %s", "Booked", "N/A"))
		lines = append(lines, fmt.Sprintf("%-22s: %s", "Injured", "N/A"))
		lines = append(lines, strings.Repeat("-", 91))
	}
	// print match stats
	homeStats := result.HomeTeam.GetStats()
	awayStats := result.AwayTeam.GetStats()
	lines = append(lines, "\nMatch Statistics")
	lines = append(lines, strings.Repeat("-", 91))
	lines = append(lines, fmt.Sprintf("%40s   |    %s", result.HomeTeam.GetName(), result.AwayTeam.GetName()))
	lines = append(lines, fmt.Sprintf("%-22s:         %3d        |       %d", "Attempts on Goal", homeStats.ShotsOnTarget+homeStats.ShotsOffTarget, awayStats.ShotsOnTarget+awayStats.ShotsOffTarget))
	// lines = append(lines, fmt.Sprintf("%-22s:         %3d        |       %d", "Blocked Shots", 1, 1))
	lines = append(lines, fmt.Sprintf("%-22s:         %3d        |       %d", "Off target", homeStats.ShotsOffTarget, awayStats.ShotsOffTarget))
	lines = append(lines, fmt.Sprintf("%-22s:         %3d        |       %d", "On target", homeStats.ShotsOnTarget, awayStats.ShotsOnTarget))
	lines = append(lines, fmt.Sprintf("%-22s:         %3d        |       %d", "Possession", homeStats.Possession, awayStats.Possession))
	lines = append(lines, fmt.Sprintf("%-22s:         %3d        |       %d", "Corners", 1, 1))
	lines = append(lines, fmt.Sprintf("%-22s:         %3d        |       %d", "Red Cards", homeStats.RedCards, awayStats.RedCards))
	lines = append(lines, fmt.Sprintf("%-22s:         %3d        |       %d", "Yellow Cards", homeStats.YellowCards, awayStats.YellowCards))
	lines = append(lines, fmt.Sprintf("%-22s:         %3d        |       %d", "Subs Used", 1, 1))
	lines = append(lines, fmt.Sprintf("%-22s:         %3d        |       %d", "Fouls", 1, 1))
	lines = append(lines, fmt.Sprintf("%-22s:         %3d        |       %d", "Key Tackles", homeStats.Tackles, awayStats.Tackles))
	lines = append(lines, fmt.Sprintf("%-22s:         %3d        |       %d", "Key Passes", homeStats.Passes, awayStats.Passes))
	lines = append(lines, strings.Repeat("-", 91))
	lines = append(lines, fmt.Sprintf("%-22s: %s %d - %d %s", "Final Score", result.HomeTeam.GetName(), len(homeStats.Goals), len(awayStats.Goals), result.AwayTeam.GetName()))
	lines = append(lines, strings.Repeat("-", 91))
	// print player stats (home)
	for i, t := range teams {
		lines = append(lines, "\nPlayer statistics - "+t.GetName())
		lines = append(lines, "Name          Pos  St  Tk  Ps  Sh  Ag | Min Sav Ktk Kps Ass Sht Gls Yel Red KAb TAb PAb SAb")
		lines = append(lines, strings.Repeat("-", 91))
		for _, p := range t.GetLineup() {
			stats := p.GetStats()
			lines = append(lines, fmt.Sprintf("%-13s %3s %3d %3d %3d %3d %3d | %3d %3d %3d %3d %3d %3d %3d %3d %3d %3d %3d %3d %3d",
				p.GetName(),
				p.GetPosition(),
				p.GetBaseAbility().Goalkeeping,
				p.GetBaseAbility().Tackling,
				p.GetBaseAbility().Passing,
				p.GetBaseAbility().Shooting,
				p.GetBaseAbility().Aggression,
				stats.MinutesPlayed,
				stats.Saves,
				stats.KeyTackles,
				stats.KeyPasses,
				stats.Assists,
				stats.ShotsOnTarget+stats.ShotsOffTarget,
				len(stats.Goals),
				utils.BoolToInt(stats.IsCautioned),
				utils.BoolToInt(stats.IsSentOff),
				p.GetMatchAbility().GoalkeepingAbs,
				p.GetMatchAbility().TacklingAbs,
				p.GetMatchAbility().PassingAbs,
				p.GetMatchAbility().ShootingAbs,
			))
		}
		if i == 0 {
			lines = append(lines, strings.Repeat("-", 91))
		}
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

// saveAsYaml saves a MatchResult to a file in YAML format.
//
// NOTE: This is not yet implemented.
func (mr *MatchResultFileStore) saveAsYaml(result *models.MatchResult, commentary []string, options MatchResultFileStoreOptions) error {
	return errors.ErrUnsupported
}

func (mr *MatchResultFileStore) saveAsJson(result *models.MatchResult, commentary []string, options MatchResultFileStoreOptions) error {
	return errors.ErrUnsupported
}

func (mr *MatchResultFileStore) Save(result *models.MatchResult, commentary []string, options MatchResultFileStoreOptions) error {
	ext := filepath.Ext(options.OutputFile)
	switch ext {
	case ".txt":
		return mr.saveAsText(result, commentary, options)
	case ".yaml":
		return mr.saveAsYaml(result, commentary, options)
	case ".json":
		return mr.saveAsJson(result, commentary, options)
	default:
		return fmt.Errorf("unsupported file format: %s", ext)
	}
}
