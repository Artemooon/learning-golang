// Package records contains the application's record data and storage.
package records

import (
	"errors"
)

var ErrNotFound = errors.New("record not found")

// Record is one vintage jazz release in the catalogue.
type Record struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Year   int    `json:"year"`
	Label  string `json:"label"`
}

// Store keeps records in memory. Replace this package with a database
// implementation when the API needs persistence.
type Store struct {
	records []Record
}

func NewStore(initial []Record) *Store {
	return &Store{records: append([]Record(nil), initial...)}
}

func (s *Store) All() []Record {
	return append([]Record(nil), s.records...)
}

func (s *Store) Add(record Record) error {
	for _, existing := range s.records {
		if existing.ID == record.ID {
			return errors.New("record ID already exists")
		}
	}
	s.records = append(s.records, record)
	return nil
}

func (s *Store) ByID(id string) (Record, error) {
	for _, record := range s.records {
		if record.ID == id {
			return record, nil
		}
	}
	return Record{}, ErrNotFound
}
