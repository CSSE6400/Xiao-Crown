package util

import (
	"encoding/csv"
	"os"

	"github.com/go-playground/log/v7"
)

// ReadFile reads csv file
func ReadFile(filename string) ([][]string, error) {
	csvFile, err := os.Open(filename)
	if err != nil {
		log.Error("read file error: %v", err)
		return nil, err
	}
	defer csvFile.Close()
	ReadCsv := csv.NewReader(csvFile)

	stringValue, _ := ReadCsv.ReadAll()
	return stringValue, nil
}
