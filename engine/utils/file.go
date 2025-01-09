package utils

import (
	"bufio"
	"os"
)

func ReadFile(filePath string, rowCallback func(string, int) error) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	rowsRead := 0
	for scanner.Scan() {
		rowsRead++
		if rowCallback != nil {
			err := rowCallback(scanner.Text(), rowsRead)
			if err != nil {
				return rowsRead, err
			}
		}
	}

	return rowsRead, scanner.Err()
}
