package store

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateSignature(t *testing.T) {
	t.Skip("Tested in record_test.go")
}

func TestSignature_SignWithKey(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	MustCreateSchema(db)

	pkeyID, err := NewPrivateKey(db).CreatePrivateKey()
	require.NoError(t, err)
	pkey, err := NewPrivateKey(db).FindPrivateKey(pkeyID)
	require.NoError(t, err)

	const recordData = "data"
	recordID, err := NewRecord(db).CreateRecord(recordData)
	require.NoError(t, err)

	svc := NewSignature(db)
	sigID, err := svc.CreateSignature(recordID, pkeyID)
	require.NoError(t, err)

	err = svc.SignWithKey(pkey, sigID)
	require.NoError(t, err)

	var got []byte
	err = db.QueryRow(`SELECT value FROM signature WHERE id = ?`, sigID).Scan(&got)
	require.NoError(t, err)

	require.NotEmpty(t, string(got))

	msgHash := sha256.Sum256([]byte(recordData))

	err = rsa.VerifyPSS(&pkey.PublicKey, crypto.SHA256, msgHash[:], got, nil)
	require.NoError(t, err)
}

func TestSignature_EmptySignature(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	MustCreateSchema(db)

	pkeyID, err := NewPrivateKey(db).CreatePrivateKey()
	require.NoError(t, err)

	recordID, err := NewRecord(db).CreateRecord("data")
	require.NoError(t, err)

	svc := NewSignature(db)
	sigID, err := svc.CreateSignature(recordID, pkeyID)
	require.NoError(t, err)

	got, err := svc.EmptySignature()
	require.NoError(t, err)

	require.Equal(t, sigID, got.ID)
	require.NotZero(t, got.RecordID)
	require.NotZero(t, got.PrivKeyID)
}
