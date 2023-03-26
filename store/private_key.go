package store

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
)

const numBits = 2048

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
func CreatePrivateKey(db *sql.DB) (int64, error) {
	pkey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return 0, err
	}
	result, err := db.Exec(`INSERT INTO private_key (key) VALUES (?) RETURNING id`, privateKeyToPEM(pkey))
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// FindPrivateKey returns the private key with the given ID.
func FindPrivateKey(db *sql.DB, id int64) (*rsa.PrivateKey, error) {
	var key string
	err := db.QueryRow(`SELECT key FROM private_key WHERE id = ?`, id).Scan(&key)
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
