package utils

import (
	"os"
	"path/filepath"
)

func EnsureDirectoryExists(filePath string) error {
	directory := filepath.Dir(filePath)

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err = os.MkdirAll(directory, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
