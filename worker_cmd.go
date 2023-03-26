package main

import "github.com/spf13/cobra"

func WorkerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "worker",
		Short: "Worker signs records in the background",
	}

	return cmd
}
