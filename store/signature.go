package store

import "database/sql"

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
