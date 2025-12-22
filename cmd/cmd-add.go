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
	Use:   "add [file]",
	Short: "Add transactions from CSV/OFX(XML FORMAT) to the database",
	Long: `The add command imports transactions from a specified CSV/OFX(XML FORMAT) file and adds them to the database.
	After adding, it displays all transactions currently stored in the database.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		ext := strings.ToLower(filepath.Ext(filePath))

		var transactions []*models.Transaction
		var err error

		switch ext {
		case ".csv":
			{
				transactions, err = services.ImportCSV(filePath)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
			}
		case ".ofx":
			{
				transactions, err = services.ImportOFX(filePath)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
			}
		default:
			fmt.Println("Unsupported file format. Please provide a CSV or OFX file.")
			return
		}

		err = services.AddTransactions(transactions)
		if err != nil {
			fmt.Println("Error saving transactions:", err)
			return
		}

		fmt.Printf("Added %d transactions to the database successfully!\n", len(transactions))
		services.PrintDBTransactions()

		// check for budget alerts after we add transactions from csv/ofx
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
