package services

import (
	"database/sql"
)

// alert structure for budget checking, it can either be overflowed(>100%) or a warning is also displayed at 90%
type BudgetAlert struct {
	Category string
	Spent    float64
	Limit    float64
	Percent  float64
	Level    string
}

// function to check the budget alerts
func CheckBudgetAlerts(db *sql.DB) ([]BudgetAlert, error) {
	rows, err := db.Query(`SELECT category, budget_limit FROM budgets`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []BudgetAlert

	for rows.Next() {
		var category string
		var limit float64

		if err := rows.Scan(&category, &limit); err != nil {
			return nil, err
		}

		spent, err := GetSpentByCategory(db, category)
		if err != nil {
			return nil, err
		}

		percent := (spent / limit) * 100
		level := "OK"

		if percent >= 100 {
			level = "ALERT!!!!!!!!!!!"
		} else if percent >= 90 {
			level = "ALMOST MAX BUDGET REACHED"
		}

		alerts = append(alerts, BudgetAlert{
			Category: category,
			Spent:    spent,
			Limit:    limit,
			Percent:  percent,
			Level:    level,
		})
	}

	return alerts, nil
}
