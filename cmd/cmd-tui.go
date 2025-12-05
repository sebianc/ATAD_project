package cmd

import (
	"fmt"

	"atad_project/services"

	"github.com/spf13/cobra"
)

var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "Browse transactions in an interactive TUI",
	RunE: func(cmd *cobra.Command, args []string) error {

		transactions, err := services.GetAllTransactions()
		if err != nil {
			return fmt.Errorf("error loading transactions: %w", err)
		}

		return services.RunTransactionTUI(transactions)
	},
}

func init() {
	mainCmd.AddCommand(browseCmd)
}
