package services

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"os"
	"sync"

	"atad_project/models"

	"github.com/olekukonko/tablewriter"

	_ "modernc.org/sqlite"
)

var db *sql.DB
var once sync.Once
var db_path string = "db/atad_project.db"

// TODO: function to initialize the database if not exists
// TODO: error handling better maybe

// function to get db instance
func InitDB() *sql.DB {
	once.Do(func() {
		var err error
		db, err = sql.Open("sqlite", db_path)
		if err != nil {
			log.Fatal(err)
		}
	})
	return db
}

// function to detect category
func DetectCategory(desc string) string {
	for category, regex := range models.CategoryRules {
		if regex.MatchString(desc) {
			return category
		}
	}
	return "Other"
}

// function used to add multiple transactions to db
func AddTransactions(transactions []*models.Transaction) error {
	db := InitDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO transactions (date, amount, description, category) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, t := range transactions {
		t.CATEGORY = DetectCategory(t.DESCRIPTION)
		_, err := stmt.Exec(t.DATE, t.AMOUNT, t.DESCRIPTION, t.CATEGORY)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// function to pretty print the content of the db
func PrintDBTransactions() ([]*models.Transaction, error) {
	db := InitDB()
	rows, err := db.Query("SELECT date, amount, description, category FROM transactions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var transactions []*models.Transaction
	// prepare table writer
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Date", "Amount", "Description", "Category"})

	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(&t.DATE, &t.AMOUNT, &t.DESCRIPTION, &t.CATEGORY)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &t)
		// append row to table
		table.Append([]string{
			t.DATE,
			fmt.Sprintf("%.2f", t.AMOUNT),
			t.DESCRIPTION,
			t.CATEGORY,
		})
	}

	// render table to stdout
	table.Render()

	return transactions, nil
}

// function to get all transactions from db in a list format
func GetAllTransactions() ([]*models.Transaction, error) {
	db := InitDB()
	rows, err := db.Query("SELECT date, amount, description, category FROM transactions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var transactions []*models.Transaction
	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(&t.DATE, &t.AMOUNT, &t.DESCRIPTION, &t.CATEGORY)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &t)
	}
	return transactions, nil
}

// function to set budget
func SetBudget(db *sql.DB, category string, limit float64) error {
	_, err := db.Exec(`
    INSERT INTO budgets (category, budget_limit)
    VALUES (?, ?)
    ON CONFLICT(category) DO UPDATE SET budget_limit = excluded.budget_limit;
`, category, limit)

	return err
}

// function to get spendings per category
func GetSpentByCategory(db *sql.DB, category string) (float64, error) {
	var spent float64
	err := db.QueryRow(`
        SELECT COALESCE(SUM(AMOUNT), 0) 
        FROM transactions 
        WHERE CATEGORY = ? AND AMOUNT < 0
    `, category).Scan(&spent)
	return math.Abs(spent), err
}

// function to get budget limit per category
func GetBudget(db *sql.DB, category string) (float64, error) {
	var limit float64
	err := db.QueryRow(`
        SELECT budget_limit FROM budgets WHERE category = ?
    `, category).Scan(&limit)
	return limit, err
}
