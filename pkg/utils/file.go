package utils

import (
	"bufio"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

func FileExists(filePath string) bool {
	st, err := os.Stat(filePath)
	return err == nil && !st.IsDir()
}

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

func FindDirEntry(dir string, compareFunc func(os.DirEntry) bool) (os.DirEntry, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if compareFunc(e) {
			return e, nil
		}
	}

	return nil, nil
}

// Traverse upwards from the current directory and print each parent directory
func FindAncestorDir(startDir string, targetDirName string, traverse bool) (string, error) {
	currentDir := startDir

	zap.L().Debug("Checking directory", zap.String("dir", currentDir), zap.String("target", targetDirName))
	entry, err := FindDirEntry(currentDir, func(e os.DirEntry) bool {
		return e.IsDir() && filepath.Base(e.Name()) == targetDirName
	})
	if err != nil {
		return "", err
	} else if entry == nil {
		parentDir := filepath.Dir(currentDir)
		if parentDir != currentDir && traverse {
			return FindAncestorDir(parentDir, targetDirName, traverse)
		} else {
			return "", nil
		}
	}

	return filepath.Join(currentDir, entry.Name()), nil
}
