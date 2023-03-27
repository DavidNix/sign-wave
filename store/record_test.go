package store

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindAvailableRecords(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		db, err := sql.Open("sqlite", ":memory:")
		require.NoError(t, err)
		MustCreateSchema(db)

		svc := NewRecord(db)

		r1, err := svc.CreateRecord("foo")
		require.NoError(t, err)
		pkey1, err := NewPrivateKey(db).CreatePrivateKey()
		require.NoError(t, err)
		_, err = NewSignature(db).CreateSignature(r1, pkey1)
		require.NoError(t, err)

		r2, err := svc.CreateRecord("bar")
		require.NoError(t, err)

		got, err := svc.FindAvailableRecords(100)
		require.NoError(t, err)

		require.Equal(t, []int64{r2}, got)
	})

	t.Run("no records", func(t *testing.T) {
		db, err := sql.Open("sqlite", ":memory:")
		require.NoError(t, err)
		MustCreateSchema(db)

		svc := NewRecord(db)

		_, err = svc.FindAvailableRecords(100)
		require.Error(t, err)
		require.EqualError(t, err, "no records found")
	})
}
