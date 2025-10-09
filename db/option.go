// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// Option represents a key-value option in the database.
type Option struct {
	ID        int64
	Key       string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// OptionRepository handles database operations for options.
type OptionRepository struct {
	db     *sql.DB
	driver string
}

// NewOptionRepository creates a new option repository.
func NewOptionRepository(db *sql.DB) *OptionRepository {
	return &OptionRepository{
		db:     db,
		driver: GetDriver(),
	}
}

// Create inserts a new option into the database.
func (r *OptionRepository) Create(key, value string) error {
	var query string
	if r.driver == "postgres" || r.driver == "postgresql" {
		query = "INSERT INTO options (key, value) VALUES ($1, $2)"
	} else {
		query = "INSERT INTO options (key, value) VALUES (?, ?)"
	}
	_, err := r.db.Exec(query, key, value)
	return err
}

// Get retrieves an option by key.
func (r *OptionRepository) Get(key string) (*Option, error) {
	option := &Option{}
	var query string
	if r.driver == "postgres" || r.driver == "postgresql" {
		query = `SELECT id, key, value, created_at, updated_at
		FROM options
		WHERE key = $1`
	} else {
		query = `SELECT id, key, value, created_at, updated_at
		FROM options
		WHERE key = ?`
	}
	err := r.db.QueryRow(query, key).Scan(
		&option.ID,
		&option.Key,
		&option.Value,
		&option.CreatedAt,
		&option.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return option, nil
}

// Update updates an option value.
func (r *OptionRepository) Update(key, value string) error {
	var query string
	if r.driver == "postgres" || r.driver == "postgresql" {
		query = `UPDATE options SET
			value = $1, updated_at = $2
		WHERE key = $3`
	} else {
		query = `UPDATE options SET
			value = ?, updated_at = ?
		WHERE key = ?`
	}
	_, err := r.db.Exec(query, value, time.Now().UTC(), key)
	return err
}

// Delete removes an option from the database.
func (r *OptionRepository) Delete(key string) error {
	var query string
	if r.driver == "postgres" || r.driver == "postgresql" {
		query = "DELETE FROM options WHERE key = $1"
	} else {
		query = "DELETE FROM options WHERE key = ?"
	}
	_, err := r.db.Exec(query, key)
	return err
}

// List retrieves all options from the database.
func (r *OptionRepository) List() ([]*Option, error) {
	rows, err := r.db.Query("SELECT id, key, value, created_at, updated_at FROM options ORDER BY key")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []*Option
	for rows.Next() {
		option := &Option{}
		if err := rows.Scan(
			&option.ID,
			&option.Key,
			&option.Value,
			&option.CreatedAt,
			&option.UpdatedAt,
		); err != nil {
			return nil, err
		}
		options = append(options, option)
	}

	return options, rows.Err()
}
