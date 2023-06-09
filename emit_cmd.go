package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DavidNix/sign-wave/emit"
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
	logger := log.New(os.Stderr, "[EMIT] ", log.LstdFlags|log.Lmsgprefix)

	db, err := openDB(cmd)
	if err != nil {
		return err
	}
	defer db.Close()

	size, err := cmd.Flags().GetInt("batch")
	if err != nil {
		return err
	}
	baseURL, err := cmd.Flags().GetString("ingest-url")
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 60 * time.Second}

	for {
		func() {
			found, err := store.NewRecord(db).FindAvailableRecords(size)
			if err != nil {
				logger.Println("Error finding records:", err)
				time.Sleep(time.Second)
				return
			}
			logger.Printf("Sending %d records\n", len(found))

			if err := emit.SendToIngest(client, baseURL, found); err != nil {
				logger.Println("Error sending records:", err)
				return
			}
			logger.Printf("Successfully sent %d records\n", len(found))
			time.Sleep(100 * time.Millisecond)
		}()
	}
}
