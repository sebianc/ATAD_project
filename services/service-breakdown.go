package services

import (
	"database/sql"
	"fmt"
	"strings"
)

// function to get category breakdown for a given month and year in ASCII format
func GetCategoryBreakdownASCII(db *sql.DB, year int, month int) (map[string]float64, error) {
	query := `
		SELECT CATEGORY, SUM(AMOUNT)
		FROM transactions
		WHERE AMOUNT < 0
		  AND substr(DATE, 1, 7) = ?
		GROUP BY CATEGORY
	`

	rows, err := db.Query(query, fmt.Sprintf("%04d-%02d", year, month))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[string]float64{}

	for rows.Next() {
		var cat string
		var amt float64
		rows.Scan(&cat, &amt)
		result[cat] = -amt
	}

	return result, nil
}

// function to print category breakdown as a bar chart in ASCII
func PrintCategoryBarChart(data map[string]float64) {
	maxWidth := 40

	for cat, value := range data {
		bars := int(value / 10) // 1 bar = 10 units
		if bars > maxWidth {
			bars = maxWidth
		}
		fmt.Printf("%-15s | %s %.2f\n", cat, strings.Repeat("█", bars), value)
	}
}
