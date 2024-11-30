package main

import (
	"html/template"
	"net/http"
	"strings"
)

type SublicenseDB struct {
	Id            int
	Name          string
	NumberOfSeats int
	LicenseKey    string
	ExpiryDate    int64
	Activ         bool
}

type Sublicense struct {
	Id            int    `json:"-"`
	Name          string `json:"name"`
	NumberOfSeats int    `json:"number_of_seats"`
	LicenseKey    string `json:"license_key"`
	ExpiryDate    string `json:"expiry_date"`
	Activ         bool   `json:"-"`
	EditLink      string `json:"-"`
	DeleteLink    string `json:"-"`
}

type SublicenseScreen struct {
	Alias fileAlias
	Data  Sublicense
}

type APIServer struct {
	listenAddr string
	db         Storage
	views      Views
}

type apiFunc func(http.ResponseWriter, *http.Request)

type Views struct {
	pages []*template.Template
}

type page int

const (
	HOMEPAGE   page = 0
	SUBLICENSE page = 1
	LOGIN      page = 2
)

type LicenseKey string

func (l *LicenseKey) valid() bool {
	fields := strings.Split(string(*l), "-")
	if len(fields) == 4 && len(fields[0]) == 4 && len(fields[1]) == 4 && len(fields[2]) == 4 && len(fields[3]) == 4 {
		return true
	}
	return false
}
