package store

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"database/sql"
)

type Signature struct {
	db *sql.DB
}

func NewSignature(db *sql.DB) *Signature {
	return &Signature{db: db}
}

// CreateSignature creates a new signature for the given record and private key.
// It does not set the value of the signature, which must be updated later.
func (svc Signature) CreateSignature(recordID, privKeyID int64) (int64, error) {
	result, err := svc.db.Exec(`
	INSERT INTO signature (record_id, private_key_id) VALUES (?, ?)
	`, recordID, privKeyID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// SignWithKey signs the associated record with the given signature ID using the given private key.
func (svc Signature) SignWithKey(pkey *rsa.PrivateKey, sigID int64) error {
	var data []byte
	err := svc.db.QueryRow(`SELECT data FROM record 
	INNER JOIN signature ON signature.record_id = record.id
	WHERE signature.id = ? LIMIT 1`, sigID).Scan(&data)
	if err != nil {
		return err
	}

	hash := sha256.Sum256(data)
	sig, err := rsa.SignPSS(rand.Reader, pkey, crypto.SHA256, hash[:], nil)
	if err != nil {
		return err
	}

	_, err = svc.db.Exec(`UPDATE signature SET value = ? WHERE id = ?`, sig, sigID)
	return err
}

type SignatureRow struct {
	ID        int64
	RecordID  int64
	PrivKeyID int64
}

// EmptySignature returns the first signature that has yet to be signed.
func (svc Signature) EmptySignature() (SignatureRow, error) {
	var row SignatureRow
	err := svc.db.QueryRow(`SELECT id, record_id, private_key_id FROM signature WHERE value IS NULL LIMIT 1`).
		Scan(&row.ID, &row.RecordID, &row.PrivKeyID)
	return row, err
}
