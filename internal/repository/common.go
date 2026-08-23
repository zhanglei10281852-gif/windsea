package repository

import (
	"database/sql"
	"fmt"
	"time"
)

type Repository struct{ DB *sql.DB }

func New(db *sql.DB) *Repository  { return &Repository{DB: db} }
func text(value time.Time) string { return value.UTC().Format(time.RFC3339Nano) }
func parse(value string) (time.Time, error) {
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse timestamp: %w", err)
	}
	return parsed, nil
}
func optional(value sql.NullString) (*time.Time, error) {
	if !value.Valid || value.String == "" {
		return nil, nil
	}
	parsed, err := parse(value.String)
	return &parsed, err
}
