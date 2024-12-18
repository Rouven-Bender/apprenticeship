package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	GetAllSublicenses() ([]*Sublicense, error)
	GetAllActivSublicenses() ([]*Sublicense, error)
	GetSublicense(id int) (*Sublicense, error)
	GetHashForUser(username string) ([]byte, error)
	CreateSublicense(lic *Sublicense) error
	CreateLoginCredentials(username string, password []byte) error
	DeleteSublicenseById(id int) error
	UpdateSublicense(lic *Sublicense) error
	DeactivateExpiredLicenses() error

	convertFromDBRepresentation(lic *SublicenseDB) *Sublicense
	convertToDBRepresentation(lic *Sublicense) (*SublicenseDB, error)
	scanIntoSublicense(rows *sql.Rows) (*SublicenseDB, error)
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

func (s *sqliteStore) GetHashForUser(username string) ([]byte, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return []byte{0}, err
	}
	defer tx.Commit()
	query := `select * from login_cred where username = ?`
	row, err := tx.Query(query, username)
	if err != nil {
		tx.Rollback()
		return []byte{0}, err
	}
	user := struct {
		username string
		pwd_hash []byte
	}{}
	for row.Next() {
		row.Scan(
			&user.username,
			&user.pwd_hash,
		)
		return user.pwd_hash, nil
	}
	return []byte{0}, fmt.Errorf("Didn't find the user")
}

func (s *sqliteStore) CreateLoginCredentials(username string, password []byte) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()
	query := `insert into login_cred (username, pwd_hash) values (?, ?)`
	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec(username, password)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (s *sqliteStore) GetSublicense(id int) (*Sublicense, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Commit()
	query := `select * from sublicenses where id = ?`
	rows, err := tx.Query(query, id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if rows.Next() {
		lic, err := s.scanIntoSublicense(rows)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		return s.convertFromDBRepresentation(lic), nil
	}
	return nil, fmt.Errorf("sublicense %d not found", id)
}

func (s *sqliteStore) GetAllSublicenses() ([]*Sublicense, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Commit()
	query := "select * from sublicenses"
	rows, err := tx.Query(query)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	slics := []*Sublicense{}
	for rows.Next() {
		lic, err := s.scanIntoSublicense(rows)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		clic := s.convertFromDBRepresentation(lic)
		slics = append(slics, clic)
	}
	return slics, nil
}

func (s *sqliteStore) GetAllActivSublicenses() ([]*Sublicense, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Commit()
	query := "select * from sublicenses where activ=true"
	rows, err := tx.Query(query)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	slics := []*Sublicense{}
	for rows.Next() {
		lic, err := s.scanIntoSublicense(rows)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		clic := s.convertFromDBRepresentation(lic)
		slics = append(slics, clic)
	}
	return slics, nil
}

func (s *sqliteStore) CreateSublicense(lic *Sublicense) error {
	clic, err := s.convertToDBRepresentation(lic)
	if err != nil {
		return err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()
	query := `insert into 
	sublicenses (name, numberOfSeats, licenseKey, expiryDate, activ)
	values (?, ?, ?, ?, ?)`
	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(clic.Name, clic.NumberOfSeats, clic.LicenseKey, clic.ExpiryDate, clic.Activ)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (s *sqliteStore) DeleteSublicenseById(id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()
	query := `delete from sublicenses where id=?`
	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (s *sqliteStore) UpdateSublicense(lic *Sublicense) error {
	clic, err := s.convertToDBRepresentation(lic)
	if err != nil {
		return err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()
	query := `update sublicenses set
		name = ?,
		numberOfSeats = ?,
		licenseKey = ?,
		expiryDate = ?,
		activ = ?
	where id = ? `
	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(clic.Name, clic.NumberOfSeats, clic.LicenseKey, clic.ExpiryDate, clic.Activ, clic.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (s *sqliteStore) DeactivateExpiredLicenses() error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()
	query := `update sublicenses set
		activ = 0
	where expiryDate < ? `
	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(time.Now().Unix())
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (s *sqliteStore) convertFromDBRepresentation(lic *SublicenseDB) *Sublicense {
	return &Sublicense{
		Id:            lic.Id,
		Name:          lic.Name,
		NumberOfSeats: lic.NumberOfSeats,
		LicenseKey:    lic.LicenseKey,
		ExpiryDate:    UnixtimeToHTMLDateString(lic.ExpiryDate),
		Activ:         lic.Activ,
		EditLink:      fmt.Sprintf("/edit/%d", lic.Id),
		DeleteLink:    fmt.Sprintf("/delete/%d", lic.Id),
	}
}

func (s *sqliteStore) convertToDBRepresentation(lic *Sublicense) (*SublicenseDB, error) {
	date, err := HTMLDateStringToUnixtime(lic.ExpiryDate)
	if err != nil {
		return nil, err
	}
	return &SublicenseDB{
		Id:            lic.Id,
		Name:          lic.Name,
		NumberOfSeats: lic.NumberOfSeats,
		LicenseKey:    lic.LicenseKey,
		ExpiryDate:    date,
		Activ:         lic.Activ,
	}, nil
}

func (s *sqliteStore) scanIntoSublicense(rows *sql.Rows) (*SublicenseDB, error) {
	lic := new(SublicenseDB)
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
