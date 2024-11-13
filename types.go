package main

import (
	"html/template"
	"net/http"
	"strings"
)

type Sublicense struct {
	Id            int
	Name          string
	NumberOfSeats int
	LicenseKey    string
	ExpiryDate    int64
	Activ         bool
	Link          string
}

type SublicenseScreen struct {
	Alias       fileAlias
	Data        Sublicense
	Conversions SublicenseConverions
}

type SublicenseConverions struct {
	Date string
}

type APIServer struct {
	listenAddr string
	db         sqliteStore
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
)

type LicenseKey string

func (l *LicenseKey) valid() bool {
	fields := strings.Split(string(*l), "-")
	if len(fields) == 4 && len(fields[0]) == 4 && len(fields[1]) == 4 && len(fields[2]) == 4 && len(fields[3]) == 4 {
		return true
	}
	return false
}
