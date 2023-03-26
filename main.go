package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "signwave",
		Short: "signwave is a suite of tools for signing records in batch",
	}

	rootCmd.AddCommand(EmitCmd())
	rootCmd.AddCommand(IngestCmd())
	rootCmd.AddCommand(WorkerCmd())

	rootCmd.PersistentFlags().String("db", "signwave.db", "Location of sqlite database file")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
