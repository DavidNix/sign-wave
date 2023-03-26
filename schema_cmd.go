package main

import (
	"github.com/DavidNix/sign-wave/store"
	"github.com/spf13/cobra"
)

func CreateSchemaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schema",
		Short: "Create the schema in the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := openDB(cmd)
			if err != nil {
				return err
			}
			defer db.Close()
			return store.CreateSchema(db)
		},
	}

	return cmd
}
