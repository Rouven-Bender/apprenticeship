package main

import (
	"html/template"
	"net/http"
)

type Sublicense struct {
	Id            int
	Name          string
	NumberOfSeats int
	LicenseKey    string
	ExpiryDate    int64
	Activ         bool
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

type Views struct {
	pages []*template.Template
}

type page int

const (
	HOMEPAGE   page = 0
	SUBLICENSE page = 1
)

type apiFunc func(http.ResponseWriter, *http.Request)
