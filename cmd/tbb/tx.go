package main

import (
	"fmt"
	"os"
	"tbb_blockchain/database"

	"github.com/spf13/cobra"
)

func txAddCmd() *cobra.Command {
	var flagFrom string = "from"
	var flagTo string = "to"
	var flagValue string = "value"
	// var flagReward string = "reward"
	var txAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Add new TX to database.",
		Run: func(cmd *cobra.Command, args []string) {
			from, _ := cmd.Flags().GetString(flagFrom)
			to, _ := cmd.Flags().GetString(flagTo)
			value, _ := cmd.Flags().GetUint(flagValue)
			// reward, _ := cmd.Flags().GetString(flagReward)

			fromAcc := database.NewAccount(from)
			toAcc := database.NewAccount(to)

			tx := database.NewTx(fromAcc, toAcc, value, "")

			state, _, err := database.NewStateFromDisk()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			defer state.Close()

			err = state.Add(tx)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			snapshot, err := state.Persist()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Printf("New DB Snapshot %x\n", snapshot)
			fmt.Println("TX successfully added to the ledger")
		},
	}
	txAddCmd.Flags().String(flagFrom, "", "From what account to send tokens")
	txAddCmd.MarkFlagRequired(flagFrom)

	txAddCmd.Flags().String(flagTo, "", "To what account to send tokens")
	txAddCmd.MarkFlagRequired(flagTo)

	txAddCmd.Flags().Uint(flagValue, 0, "How many tokens to send")
	txAddCmd.MarkFlagRequired(flagValue)
	return txAddCmd
}

func txCmd() *cobra.Command {
	var txsCmd = &cobra.Command{
		Use:   "tx",
		Short: "Interact with txs (add...).",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageErr()
		},
		Run: func(cmd *cobra.Command, args []string) {},
	}
	txsCmd.AddCommand(txAddCmd())
	return txsCmd
}
