package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	scanIntoSublicenses(*sql.Rows) (*Sublicense, error)
	GetAllSublicenses() ([]*Sublicense, error)
	GetSublicense(id int) (*Sublicense, error)
}

type sqliteStore struct {
	db *sql.DB
}

func NewSqliteStore() (*sqliteStore, error) {
	db, err := sql.Open("sqlite3", "db.db")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &sqliteStore{
		db: db,
	}, nil
}

func (s *sqliteStore) GetSublicense(id int) (*Sublicense, error) {
	rows, err := s.db.Query("select * from sublicenses where id = ?", id)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return s.scanIntoSublicense(rows)
	}
	return nil, fmt.Errorf("sublicense %d not found", id)
}

func (s *sqliteStore) GetAllSublicenses() ([]*Sublicense, error) {
	rows, err := s.db.Query("select * from sublicenses")
	if err != nil {
		return nil, err
	}
	slics := []*Sublicense{}
	for rows.Next() {
		lics, err := s.scanIntoSublicense(rows)
		if err != nil {
			return nil, err
		}
		slics = append(slics, lics)
	}
	return slics, nil
}

func (s *sqliteStore) scanIntoSublicense(rows *sql.Rows) (*Sublicense, error) {
	lic := new(Sublicense)
	err := rows.Scan(
		&lic.Id,
		&lic.Name,
		&lic.NumberOfSeats,
		&lic.LicenseKey,
		&lic.ExpiryDate,
		&lic.Activ,
	)
	if err != nil {
		return nil, err
	}
	return lic, nil
}
