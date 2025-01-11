package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/esmshub/esms-go/engine/models"
)

var MIN_ROSTER_ATTR_ACOUNT = 18
var MAX_ROSTER_ATTR_COUNT = 30

func readFields(row string, min, max int) ([]string, error) {
	values := strings.Fields(row)
	valuesCount := len(values)
	if valuesCount < min || valuesCount > max {
		return nil, fmt.Errorf("got %d attr names, expected between %d-%d", valuesCount, min, max)
	}

	return values, nil
}

func atoi(value string) int {
	result, err := strconv.Atoi(value)
	if err != nil {
		return -1
	}

	return result
}

func findIndex(arr []string, target string) int {
	for i, v := range arr {
		if strings.EqualFold(v, target) {
			return i
		}
	}
	return -1
}

func LoadRoster(filePath string) ([]*models.Player, error) {
	roster := []*models.Player{}

	file, err := os.Open(filePath)
	if err != nil {
		return roster, errors.New("unable to load roster file")
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// parse header
	if !scanner.Scan() {
		return roster, errors.New("roster does not look valid")
	}

	attrNames, err := readFields(scanner.Text(), MIN_ROSTER_ATTR_ACOUNT, MAX_ROSTER_ATTR_COUNT)
	if err != nil {
		fmt.Println("Error:", err)
		return roster, errors.New("invalid roster file header")
	}

	for scanner.Scan() {
		row := scanner.Text()
		// skip separator row if present
		if strings.HasPrefix(row, "-----") {
			continue
		}

		values, err := readFields(row, MIN_ROSTER_ATTR_ACOUNT, MAX_ROSTER_ATTR_COUNT)
		if err != nil {
			fmt.Println("Error:", err)
			return roster, errors.New("invalid roster file header")
		}

		player := &models.Player{
			Name: values[findIndex(attrNames, "name")],
			Age:  atoi(values[findIndex(attrNames, "age")]),
			Nat:  values[findIndex(attrNames, "nat")],
			Abilities: &models.PlayerAbilities{
				Goalkeeping:    atoi(values[findIndex(attrNames, "st")]),
				Tackling:       atoi(values[findIndex(attrNames, "tk")]),
				Passing:        atoi(values[findIndex(attrNames, "ps")]),
				Shooting:       atoi(values[findIndex(attrNames, "sh")]),
				Aggression:     atoi(values[findIndex(attrNames, "ag")]),
				GoalkeepingAbs: atoi(values[findIndex(attrNames, "kab")]),
				TacklingAbs:    atoi(values[findIndex(attrNames, "tab")]),
				PassingAbs:     atoi(values[findIndex(attrNames, "pab")]),
				ShootingAbs:    atoi(values[findIndex(attrNames, "sab")]),
			},
			Stats: &models.PlayerStats{
				GamesStarted:       atoi(values[findIndex(attrNames, "gam")]),
				GamesSubbed:        atoi(values[findIndex(attrNames, "sub")]),
				MinutesPlayed:      atoi(values[findIndex(attrNames, "min")]),
				MomAwards:          atoi(values[findIndex(attrNames, "mom")]),
				Saves:              atoi(values[findIndex(attrNames, "sav")]),
				GoalsConceded:      atoi(values[findIndex(attrNames, "con")]),
				KeyTackles:         atoi(values[findIndex(attrNames, "ktk")]),
				KeyPasses:          atoi(values[findIndex(attrNames, "kps")]),
				Shots:              atoi(values[findIndex(attrNames, "sht")]),
				Goals:              atoi(values[findIndex(attrNames, "gls")]),
				Assists:            atoi(values[findIndex(attrNames, "ass")]),
				DisciplinaryPoints: atoi(values[findIndex(attrNames, "dp")]),
				WeeksInjured:       atoi(values[findIndex(attrNames, "inj")]),
				GamesSuspended:     atoi(values[findIndex(attrNames, "sus")]),
			},
		}
		roster = append(roster, player)
	}

	return roster, nil
}
