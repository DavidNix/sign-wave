package main

import "github.com/spf13/cobra"

const defaultIngestPort = "9876"

func IngestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ingest",
		Short: "Ingest records from an http endpoint",
	}

	cmd.Flags().String("port", defaultIngestPort, "Port to bind http server")

	return cmd
}
