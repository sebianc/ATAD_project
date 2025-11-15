package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:   "atad-cli",
	Short: "A CLI tool to track personal income and expenses",
	Long: `ATAD CLI is a command-line tool for tracking personal income and expenses. 
	Import transactions from bank statements, categorize them automatically, set budgets, and generate insightful 
	reports— all from your terminal.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to ATAD CLI! Usage: atad-cli [command] [flags]")
	},
}

func Execute() {
	cobra.CheckErr(mainCmd.Execute())
}
