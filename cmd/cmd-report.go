package cmd

import (
	"atad_project/services"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report <YYYY-MM>",
	Short: "Generate monthly spending report",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		parts := strings.Split(args[0], "-")
		year, _ := strconv.Atoi(parts[0])
		month, _ := strconv.Atoi(parts[1])

		db := services.InitDB()

		data, err := services.GetMonthlySpendingASCII(db, year, month)
		if err != nil {
			return err
		}

		services.PrintMonthlySpendingChart(data, year, month)
		return nil
	},
}

func init() {
	mainCmd.AddCommand(reportCmd)
}
