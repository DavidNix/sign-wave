package store

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func CreateSchema(db *sql.DB) error {
	_, err := db.Exec(`
  PRAGMA foreign_keys=ON;

  CREATE TABLE record (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    data TEXT NOT NULL
  );

  CREATE TABLE private_key (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key TEXT UNIQUE NOT NULL,
    leased INTEGER NOT NULL DEFAULT 0,
    used_count INTEGER DEFAULT 0
  );

  CREATE TABLE signature (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    record_id INTEGER UNIQUE,
    private_key_id INTEGER,
    value TEXT,
    FOREIGN KEY (record_id) REFERENCES record(id),
    FOREIGN KEY (private_key_id) REFERENCES private_key(id)
  );
`)
	return err
}

func MustCreateSchema(db *sql.DB) {
	if err := CreateSchema(db); err != nil {
		panic(err)
	}
}
