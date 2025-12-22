package cmd

import (
	"atad_project/models"
	"atad_project/services"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// add command definition and behavior
// it will take as an input the CSV/OFX(XML FORMAT) file, import the transactions and add them to the database. also display the content of the db in the end
var addCmd = &cobra.Command{
	Use:   "add (file/f <file> | manual/m <date amount description...>)",
	Short: "Add transactions from a file or add one manually",
	Long: `The add command adds transactions to the database in two ways:
	- From a CSV or OFX (XML) file using --file / -f
	- Manually adding a single income or expense using --manual / -m

	After adding, all transactions are displayed and budget alerts are checked.`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		flag := args[0]

		switch flag {
		case "file", "f":
			if len(args) < 2 {
				fmt.Println("Missing file path after -f / --file")
				return
			}

			filePath := args[1]
			ext := strings.ToLower(filepath.Ext(filePath))

			var transactions []*models.Transaction
			var err error

			switch ext {
			case ".csv":
				transactions, err = services.ImportCSV(filePath)
			case ".ofx":
				transactions, err = services.ImportOFX(filePath)
			default:
				fmt.Println("Unsupported file format. Please provide a CSV or OFX file.")
				return
			}

			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			err = services.AddTransactions(transactions)
			if err != nil {
				fmt.Println("Error saving transactions:", err)
				return
			}

			fmt.Printf("Added %d transactions to the database successfully!\n", len(transactions))
			services.PrintDBTransactions()

		case "manual", "m":
			if len(args) < 4 {
				fmt.Println("Manual mode requires: DATE AMOUNT DESCRIPTION")
				return
			}

			date := args[1]
			amount := args[2]
			description := args[3:]

			transactions, err := services.ParseManualTransaction(
				date,
				amount,
				"", // category (default)
				description...,
			)
			if err != nil {
				fmt.Println("Error parsing manual transaction:", err)
				return
			}

			err = services.AddTransactions(transactions)
			if err != nil {
				fmt.Println("Error saving transactions:", err)
				return
			}
			services.PrintDBTransactions()
		default:
			fmt.Println("Invalid flag. Use -f/--file or -m/--manual")
			return
		}

		// common: budget alerts
		db := services.InitDB()
		alerts, _ := services.CheckBudgetAlerts(db)
		for _, a := range alerts {
			if a.Level != "OK" {
				fmt.Printf("BUDGET ALERT %s: %.1f%% spent\n", a.Category, a.Percent)
			}
		}

	},
}

func init() {
	mainCmd.AddCommand(addCmd)
}
