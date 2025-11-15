package cmd

import (
	"atad_project/services"
	"fmt"

	"github.com/spf13/cobra"
)

// import command definition and behavior
// it will take as an input the CSV file, import  the transactions and display them in a table format
var importCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "Import and display transactions from a CSV file",
	Long:  `The import command reads transactions from a specified CSV file and displays them in a formatted table.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		transactions, err := services.ImportCSV(filePath)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Printf("Transactions from CSV (%d):\n", len(transactions))
		services.PrettyPrintTransactions(transactions)
	},
}

func init() {
	mainCmd.AddCommand(importCmd)
}
