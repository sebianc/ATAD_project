package cmd

import (
	"atad_project/services"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// budget command definition and behavior
// it will take as an input the category and limit amount to set a budget for that category. the alerts will be checked immediately and it can be either OK or >=90% warning or >=100% alert
var budgetSetCmd = &cobra.Command{
	Use:   "budget <category> <limit>",
	Short: "Set a budget limit for a specific category",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {

		category := args[0]

		limit, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}

		db := services.InitDB()
		err = services.SetBudget(db, category, limit)
		if err != nil {
			return fmt.Errorf("failed to set budget: %w", err)
		}

		fmt.Printf("Budget set for category '%s': %.2f\n", category, limit)

		// also check for alerts immediately after setting a budget
		alerts, _ := services.CheckBudgetAlerts(db)
		for _, a := range alerts {
			if a.Level != "OK" {
				fmt.Printf("BUDGET ALERT %s: %.1f%% spent\n", a.Category, a.Percent)
			}
		}
		return nil
	},
}

func init() {
	mainCmd.AddCommand(budgetSetCmd)
}
