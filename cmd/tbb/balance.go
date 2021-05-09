package main

import (
	"fmt"
	"os"
	"tbb_blockchain/database"

	"github.com/spf13/cobra"
)

func balancesCmd() *cobra.Command {
	var balancesCmd = &cobra.Command{
		Use:   "balances",
		Short: "Interact with balances (list...).",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageErr()
		},
		Run: func(cmd *cobra.Command, args []string) {},
	}

	var balancesListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all balance",
		Run: func(cmd *cobra.Command, args []string) {
			state, lastSnapshot, err := database.NewStateFromDisk(getDataDirFromCmd(cmd))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer state.Close()

			fmt.Printf("Accounts balances at %x...\n", lastSnapshot[0:7])
			fmt.Println("")
			for account, balance := range state.Balances {
				fmt.Println(fmt.Sprintf("%s: %d", account, balance))
			}
		},
	}

	balancesCmd.AddCommand(balancesListCmd)
	return balancesCmd
}
