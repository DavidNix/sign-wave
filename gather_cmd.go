package main

import "github.com/spf13/cobra"

func EmitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "emit",
		Short: "Continuously emit records in batches",
	}

	cmd.Flags().Uint("batch", 1000, "Batch size")

	return cmd
}
