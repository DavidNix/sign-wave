#!/usr/bin/env bash

set -euo pipefail

DB_FILE="${1:-signwave.db}"

echo "Seeding records..."
sqlite3 "$DB_FILE" "insert into record (data) select random() from generate_series(1, 100000);"
echo "Number of records in database: $(sqlite3 "$DB_FILE" "select count(*) from record;")"

echo "Seeding private keys..."
for i in {1..100}; do
  PEM="$(openssl genrsa -out /dev/stdout 2048)"
  sqlite3 "$DB_FILE" "insert into private_key (key) values ('$PEM');"
done
echo "Number of private keys in database: $(sqlite3 "$DB_FILE" "select count(*) from private_key;")"