package store

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
)

const numBits = 2048

type Database interface {
	Exec(query string, args ...any) (sql.Result, error)
	QueryRow(query string, args ...any) *sql.Row
}

type PrivateKey struct {
	db Database
}

func NewPrivateKey(db Database) *PrivateKey {
	return &PrivateKey{db: db}
}

func privateKeyToPEM(privkey *rsa.PrivateKey) string {
	pkeyBytes := x509.MarshalPKCS1PrivateKey(privkey)
	pkeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: pkeyBytes,
		},
	)
	return string(pkeyPem)
}

// CreatePrivateKey creates a new private key in the database and returns its ID.
func (svc PrivateKey) CreatePrivateKey() (int64, error) {
	pkey, err := rsa.GenerateKey(rand.Reader, numBits)
	if err != nil {
		return 0, err
	}
	result, err := svc.db.Exec(`INSERT INTO private_key (key) VALUES (?)`, privateKeyToPEM(pkey))
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// FindPrivateKey returns the private key with the given ID.
func (svc PrivateKey) FindPrivateKey(id int64) (*rsa.PrivateKey, error) {
	var key string
	err := svc.db.QueryRow(`SELECT key FROM private_key WHERE id = ?`, id).Scan(&key)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode([]byte(key))
	pkey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pkey, nil
}

// Lease finds an available private key and leases it for the given record IDs.
// Must be used with a transaction to prevent races.
func (svc PrivateKey) Lease(recordIDs []int64) error {
	var privKeyID int64
	err := svc.db.QueryRow(`SELECT id FROM private_key WHERE leased = 0 ORDER BY used_count ASC LIMIT 1`).Scan(&privKeyID)
	if err != nil {
		return err
	}

	_, err = svc.db.Exec(`UPDATE private_key SET leased = 1, used_count = used_count + 1 WHERE private_key.id = ?`, privKeyID)
	if err != nil {
		return err
	}

	for _, id := range recordIDs {
		_, err = svc.db.Exec(`INSERT INTO signature (record_id, private_key_id) VALUES (?, ?)`, id, privKeyID)
		if err != nil {
			return err
		}
	}

	return nil
}

// ResetLeases resets the leases on all private keys which have signed all associated records.
// TODO: Return how many leases were reset.
func (svc PrivateKey) ResetLeases() error {
	_, err := svc.db.Exec(`UPDATE private_key SET leased = 0
				FROM (
				SELECT private_key_id, count(CASE WHEN value IS NULL THEN 1 ELSE NULL END) AS cnt 
					FROM signature 
					INNER JOIN private_key ON private_key.id = signature.private_key_id 
					WHERE private_key.leased = 1 GROUP BY private_key_id
					HAVING cnt = 0
				) AS s 
                WHERE private_key.id = s.private_key_id`)
	return err
}
