package main

import "github.com/spf13/cobra"

func IngestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ingest",
		Short: "Ingest records from an http endpoint",
	}

	cmd.Flags().String("port", "9876", "Port to bind http server")

	return cmd
}
