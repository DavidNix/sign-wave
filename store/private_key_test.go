package store

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateFindPrivateKey(t *testing.T) {
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

func TestPrivateKey_Lease(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		db, err := sql.Open("sqlite", ":memory:")
		require.NoError(t, err)
		defer db.Close()
		MustCreateSchema(db)

		recordIDs := make([]int64, 5)
		for i := 0; i < 5; i++ {
			id, err := NewRecord(db).CreateRecord("data")
			require.NoError(t, err)
			recordIDs[i] = id
		}

		_, err = db.Exec(`INSERT INTO private_key (key, used_count) VALUES (?, ?)`, "used", 1)
		require.NoError(t, err)
		_, err = db.Exec(`INSERT INTO private_key (key, leased) VALUES (?, 1)`, "leased")
		require.NoError(t, err)

		tx, err := db.Begin()
		require.NoError(t, err)
		defer tx.Rollback()

		svc := NewPrivateKey(tx)

		pkeyID, err := svc.CreatePrivateKey()
		require.NoError(t, err)

		err = svc.Lease(recordIDs)
		require.NoError(t, err)

		var count int
		err = tx.QueryRow(`SELECT COUNT(*) FROM signature WHERE private_key_id = ? AND value IS NULL`, pkeyID).Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 5, count)

		var leased, usedCount int
		err = tx.QueryRow(`SELECT leased, used_count FROM private_key WHERE id = ?`, pkeyID).Scan(&leased, &usedCount)
		require.Equal(t, 1, leased)
		require.Equal(t, 1, usedCount)
	})

	t.Run("no available keys", func(t *testing.T) {
		db, err := sql.Open("sqlite", ":memory:")
		require.NoError(t, err)
		defer db.Close()
		MustCreateSchema(db)

		recordIDs := make([]int64, 3)
		for i := 0; i < 3; i++ {
			id, err := NewRecord(db).CreateRecord("data")
			require.NoError(t, err)
			recordIDs[i] = id
		}

		_, err = db.Exec(`INSERT INTO private_key (key, leased) VALUES (?, ?)`, "leased", 1)
		require.NoError(t, err)

		tx, err := db.Begin()
		require.NoError(t, err)
		defer tx.Rollback()

		err = NewPrivateKey(tx).Lease(recordIDs)

		require.Error(t, err)
		require.EqualError(t, err, "sql: no rows in result set")
	})
}
