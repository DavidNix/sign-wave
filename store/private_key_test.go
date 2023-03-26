package store

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrivateKey(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	MustCreateSchema(db)

	got, err := CreatePrivateKey(db)
	require.NoError(t, err)

	require.NotZero(t, got)

	got2, err := CreatePrivateKey(db)
	require.NoError(t, err)
	require.Greater(t, got2, got)

	gotPkey, err := FindPrivateKey(db, got)
	require.NoError(t, err)
	require.NoError(t, gotPkey.Validate())
}
