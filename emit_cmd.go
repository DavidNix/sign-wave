package main

import (
	"fmt"
	"time"

	"github.com/DavidNix/sign-wave/store"
	"github.com/spf13/cobra"
)

func EmitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "emit",
		Short: "Continuously emit records in batches",
		RunE:  doEmit,
	}

	cmd.Flags().Int("batch", 1000, "Batch size")
	cmd.Flags().String("ingest-url", "http://localhost:"+defaultIngestPort, "Ingest url")

	return cmd
}

func doEmit(cmd *cobra.Command, args []string) error {
	db, err := openDB(cmd)
	if err != nil {
		return err
	}
	defer db.Close()

	size, err := cmd.Flags().GetInt("batch")
	if err != nil {
		return err
	}

	for {
		found, err := store.NewRecord(db).FindAvailableRecords(size)
		if err != nil {
			fmt.Println("Error finding records: ", err)
		}
		fmt.Println("FOUND:")
		fmt.Println(found)
		time.Sleep(1 * time.Second)
	}
}
