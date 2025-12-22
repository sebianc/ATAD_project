package services

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/guptarohit/asciigraph"
)

// function to get monthly spending data for a given month and year in ASCII format
func GetMonthlySpendingASCII(db *sql.DB, year int, month int) ([]float64, error) {
	query := `
		SELECT DATE, AMOUNT
		FROM transactions
		WHERE AMOUNT < 0
		  AND substr(DATE, 1, 7) = ?
	`
	rows, err := db.Query(query, fmt.Sprintf("%04d-%02d", year, month))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 31 days in a month
	daily := make([]float64, 31)

	for rows.Next() {
		var date string
		var amount float64
		rows.Scan(&date, &amount)

		day, _ := strconv.Atoi(date[8:10])
		daily[day-1] += -amount // convert to positive values
	}

	return daily, nil
}

// function to print monthly spending as a line chart in ASCII using asciigraph library
func PrintMonthlySpendingChart(values []float64, year int, month int) {
	graph := asciigraph.Plot(
		values,
		asciigraph.Caption(fmt.Sprintf("Monthly Spending %04d-%02d", year, month)),
		asciigraph.Height(15),
	)
	fmt.Println(graph)
}
