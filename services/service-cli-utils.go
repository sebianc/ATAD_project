package services

import (
	"fmt"
	"os"

	"atad_project/models"

	"github.com/olekukonko/tablewriter"
)

// function to use tablewriter to print transactions in a table format
func PrettyPrintTransactions(transactions []*models.Transaction) {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Date", "Amount", "Description", "Category"})

	for _, t := range transactions {
		table.Append([]string{
			t.DATE,
			fmt.Sprintf("%.2f", t.AMOUNT),
			t.DESCRIPTION,
			t.CATEGORY,
		})
	}

	table.Render()
}
