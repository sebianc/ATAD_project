package services

import (
	"fmt"
	"os"

	"atad_project/models"

	"github.com/gocarina/gocsv"
)

// function to import CSV file
func ImportCSV(filePath string) ([]*models.Transaction, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	var transactions []*models.Transaction
	if err := gocsv.UnmarshalFile(file, &transactions); err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %w", err)
	}

	return transactions, nil
}
