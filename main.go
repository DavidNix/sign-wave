package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	_ "modernc.org/sqlite"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "signwave",
		Short: "signwave is a suite of tools for signing records in batch",
	}

	rootCmd.AddCommand(CreateSchemaCmd())
	rootCmd.AddCommand(EmitCmd())
	rootCmd.AddCommand(IngestCmd())
	rootCmd.AddCommand(WorkerCmd())

	rootCmd.PersistentFlags().String("db", "signwave.db", "Location of sqlite database file")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func openDB(cmd *cobra.Command) (*sql.DB, error) {
	dsn, err := cmd.Root().PersistentFlags().GetString("db")
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	return db, db.Ping()
}
