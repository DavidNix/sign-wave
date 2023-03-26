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
	svc := NewPrivateKey(db)

	got, err := svc.CreatePrivateKey()
	require.NoError(t, err)
	require.NotZero(t, got)

	got2, err := svc.CreatePrivateKey()
	require.NoError(t, err)
	require.Greater(t, got2, got)

	gotPkey, err := svc.FindPrivateKey(got)
	require.NoError(t, err)
	require.NoError(t, gotPkey.Validate())
}
