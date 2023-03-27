package main

import (
	"log"
	"os"
	"time"

	"github.com/DavidNix/sign-wave/store"
	"github.com/spf13/cobra"
)

func WorkerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "worker",
		Short: "Worker signs records in the background",
		RunE:  doWork,
	}

	return cmd
}

func doWork(cmd *cobra.Command, args []string) error {
	db, dberr := openDB(cmd)
	if dberr != nil {
		return dberr
	}

	logger := log.New(os.Stderr, "[WORKER] ", log.LstdFlags|log.Lmsgprefix)

	pkeySvc := store.NewPrivateKey(db)
	sigSvc := store.NewSignature(db)

	for {
		if err := pkeySvc.ResetLeases(); err != nil {
			logger.Println("Error resetting leases:", err)
		}
		found, err := sigSvc.EmptySignature()
		if err != nil {
			logger.Println("Error finding empty signature:", err)
			time.Sleep(1 * time.Second)
			continue
		}
		pkey, err := pkeySvc.FindPrivateKey(found.PrivKeyID)
		if err != nil {
			logger.Println("Error finding private key:", err)
			time.Sleep(1 * time.Second)
			continue
		}
		err = sigSvc.SignWithKey(pkey, found.ID)
		if err != nil {
			logger.Println("Error signing record:", err)
			time.Sleep(1 * time.Second)
			continue
		}
		logger.Printf("Signed record %d", found.RecordID)
	}
}
