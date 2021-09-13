package utils

import (
	"encoding/csv"
	"os"
)

// ReadCSV - Read all items from a csv item file
func ReadCSV(filePath string) [][]string {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil
	}
	r := csv.NewReader(file)
	r.LazyQuotes = true
	lines, err := r.ReadAll()
	if err != nil {
		return nil
	}

	return lines
}

// WriteCSV - write the csv item file for delete/update operations
func WriteCSV(filePath string, lines [][]string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	for _, line := range lines {
		if len(line[0]) > 0 {
			err = w.Write(line)
			if err != nil {
				return err
			}
		}
	}
	w.Flush()
	return nil
}
