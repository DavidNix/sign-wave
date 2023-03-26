package store

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateSchema(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	err = CreateSchema(db)

	require.NoError(t, err)
}

func TestMustCreateSchema(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	require.NotPanics(t, func() {
		MustCreateSchema(db)
	})
}
