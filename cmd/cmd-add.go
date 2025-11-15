package cmd

import (
	"atad_project/services"
	"fmt"

	"github.com/spf13/cobra"
)

// add command definition and behavior
// it will take as an input the CSV file, import the transactions and add them to the database. also display the content of the db in the end
var addCmd = &cobra.Command{
	Use:   "add [file]",
	Short: "Add transactions from CSV to the database",
	Long: `The add command imports transactions from a specified CSV file and adds them to the database.
	After adding, it displays all transactions currently stored in the database.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		transactions, err := services.ImportCSV(filePath)
		if err != nil {
			fmt.Println("Error reading CSV:", err)
			return
		}

		err = services.AddTransactions(transactions)
		if err != nil {
			fmt.Println("Error saving transactions:", err)
			return
		}

		fmt.Printf("Added %d transactions to the database successfully!\n", len(transactions))
		services.PrintDBTransactions()
	},
}

func init() {
	mainCmd.AddCommand(addCmd)
}
