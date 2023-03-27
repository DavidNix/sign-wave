package store

import (
	"database/sql"
	"errors"
)

type Record struct {
	db *sql.DB
}

func NewRecord(db *sql.DB) *Record {
	return &Record{db: db}
}

func (svc Record) CreateRecord(data string) (int64, error) {
	result, err := svc.db.Exec(`
	INSERT INTO record (data) VALUES (?)
	`, data)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// FindAvailableRecords returns a list of record IDs that have been marked for signatures.
func (svc Record) FindAvailableRecords(limit int) ([]int64, error) {
	rows, err := svc.db.Query(`
	SELECT id FROM record
	WHERE id NOT IN (SELECT record_id FROM signature)
	LIMIT ?
  `, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return nil, errors.New("no records found")
	}
	return ids, nil
}
