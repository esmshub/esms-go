package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/esmshub/esms-go/engine/models"
	"github.com/esmshub/esms-go/engine/types"
	"github.com/esmshub/esms-go/engine/validators"
	"github.com/esmshub/esms-go/pkg/utils"
	"golang.org/x/exp/slices"
)

const DefaultTacticsMatrixFileName = "tactics.dat"

func isValidTactic(tactic string) bool {
	return slices.ContainsFunc(types.ValidTactics, func(t string) bool {
		return strings.EqualFold(tactic, t)
	})
}

func parseTacticRow(line string, row int) (string, []float64, error) {
	t1 := strings.SplitAfterN(line, " ", 2)
	t2 := strings.Split(t1[0], ":")
	if len(t2) != 2 {
		return "", nil, fmt.Errorf("invalid format on row %d: \"%s\"", row, line)
	}
	tactic := strings.TrimSpace(t2[0])
	if strings.Contains(tactic, "vs") {
		counterTactics := strings.Split(tactic, "vs")
		if len(counterTactics) != 2 {
			return "", nil, fmt.Errorf("unrecognised tactic on row %d: \"%s\"", row, tactic)
		}

		for _, t := range counterTactics {
			if !isValidTactic(t) {
				return "", nil, fmt.Errorf("unrecognised tactic on row %d: \"%s\"", row, tactic)
			}
		}
	} else if !isValidTactic(tactic) {
		return "", nil, fmt.Errorf("unrecognised tactic on row %d: \"%s\"", row, tactic)
	}

	position := strings.TrimSpace(t2[1])
	if !validators.IsValidPosition(position) {
		return "", nil, fmt.Errorf("unrecognised position on row %d: \"%s\"", row, position)
	}

	matrixStr, err := utils.Substring(t1[1], "[", "]")
	if err != nil {
		return "", nil, err
	}
	t3 := strings.Split(matrixStr, ",")
	if len(t3) != 3 {
		return "", nil, fmt.Errorf("invalid value matrix on row %d: \"%s\"", row, t1[1])
	}

	data := []float64{}
	for i := 0; i < len(t3); i++ {
		val, err := strconv.ParseFloat(strings.TrimSpace(t3[i]), 64)
		if err != nil {
			return "", nil, err
		}
		data = append(data, val)
	}

	return fmt.Sprintf("%s_%s", tactic, position), data, nil
}

func LoadTactics(filePath string) (*models.TacticsMatrix, error) {
	matrix := models.TacticsMatrix{}
	_, err := utils.ReadFile(filePath, func(line string, row int) error {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			// log.Debug().Msg("Skipping empty line.")
			return nil
		}

		key, data, err := parseTacticRow(trimmedLine, row)
		if err != nil {
			return err
		}

		matrix[key] = data
		return nil
	})

	return &matrix, err
}
