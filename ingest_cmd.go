package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/DavidNix/sign-wave/emit"
	"github.com/DavidNix/sign-wave/store"
	"github.com/spf13/cobra"
)

const defaultIngestPort = "9876"

func IngestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ingest",
		Short: "Ingest records from an http endpoint",
		RunE:  doIngest,
	}

	cmd.Flags().String("port", defaultIngestPort, "Port to bind http server")

	return cmd
}

func doIngest(cmd *cobra.Command, args []string) error {
	logger := log.New(os.Stderr, "[INGEST] ", log.LstdFlags|log.Lmsgprefix)

	db, err := openDB(cmd)
	if err != nil {
		return err
	}

	setupIngestServer(logger, db)

	port, err := cmd.Flags().GetString("port")
	if err != nil {
		return err
	}
	logger.Println("Server listening on port:", port)
	return http.ListenAndServe(":"+port, nil)
}

// TODO: Move to public function and unit test.
func setupIngestServer(logger *log.Logger, db *sql.DB) {
	http.HandleFunc("/v1/batch", func(w http.ResponseWriter, r *http.Request) {
		var body emit.Request
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			logger.Println("Error decoding request:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Printf("Received batch of %d records\n", len(body.RecordIDs))
		tx, err := db.Begin()
		if err != nil {
			logger.Println("Error begin tx:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		if err = store.NewPrivateKey(tx).Lease(body.RecordIDs); err != nil {
			logger.Println("Error leasing private key:", err)
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		if err = tx.Commit(); err != nil {
			logger.Println("Error commit tx:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
