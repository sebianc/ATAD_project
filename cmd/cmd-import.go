package cmd

import (
	"atad_project/services"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// import command definition and behavior
// it will take as an input the CSV/OFX(XML FORMAT) file, import  the transactions and display them in a table format
var importCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "Import and display transactions from a CSV/OFX(XML FORMAT) file",
	Long:  `The import command reads transactions from a specified CSV/OFX(XML FORMAT) file and displays them in a formatted table.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		ext := strings.ToLower(filepath.Ext(filePath))

		switch ext {
		case ".csv":
			{
				transactions, err := services.ImportCSV(filePath)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
				fmt.Printf("Transactions from CSV (%d):\n", len(transactions))
				services.PrettyPrintTransactions(transactions)
				return
			}
		case ".ofx":
			{
				transactions, err := services.ImportOFX(filePath)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
				fmt.Printf("Transactions from OFX (%d):\n", len(transactions))
				services.PrettyPrintTransactions(transactions)
				return
			}
		default:
			fmt.Println("Unsupported file format. Please provide a CSV or OFX file.")
			return
		}
	},
}

func init() {
	mainCmd.AddCommand(importCmd)
}
