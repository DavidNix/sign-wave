#!/usr/bin/env bash

set -euo pipefail

SHOULD_SEED="${SEED:-true}"
DB_FILE="$1"

if [ -f "$DB_FILE" ]; then
  echo "Database file already exists: $DB_FILE"
  exit 1
fi

echo "Initializing schema..."
sqlite3 "$DB_FILE" <<EOF
  PRAGMA foreign_keys=ON;

  CREATE TABLE record (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    data TEXT NOT NULL
  );

  CREATE TABLE private_key (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key TEXT NOT NULL,
    used INTEGER DEFAULT 0
  );

  CREATE TABLE signature (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    record_id INTEGER UNIQUE,
    private_key_id INTEGER,
    value TEXT,
    FOREIGN KEY (record_id) REFERENCES record(id),
    FOREIGN KEY (private_key_id) REFERENCES private_key(id)
  );
EOF

if [ "$SHOULD_SEED" != "true" ]; then
  echo "Skipping seeding records."
  exit 0
fi

echo "Seeding records..."
sqlite3 "$DB_FILE" "insert into record (data) select random() from generate_series(1, 100000);"
echo "Number of records in database: $(sqlite3 "$DB_FILE" "select count(*) from record;")"

echo "Seeding private keys..."
for i in {1..100}; do
  PEM="$(openssl genrsa -out /dev/stdout 2048)"
  sqlite3 "$DB_FILE" "insert into private_key (key) values ('$PEM');"
done
echo "Number of private keys in database: $(sqlite3 "$DB_FILE" "select count(*) from private_key;")"