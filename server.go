package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
)

func (v *Views) render(p page, block string, status int, w http.ResponseWriter, d interface{}) error {
	w.WriteHeader(status)
	return v.pages[p].ExecuteTemplate(w, block, d)
}

func NewAPIServer(listenAddr string, database sqliteStore) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		db:         database,
		views:      loadTemplates(),
	}
}

func (s *APIServer) Run() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /cdn/{filename}", http.StripPrefix("/cdn/", fs))
	mux.Handle("GET /api/html/table", makeHTTPHandleFunc(s.table))
	mux.Handle("GET /edit/{id}", makeHTTPHandleFunc(s.edit))
	mux.Handle("PATCH /edit/{id}", makeHTTPHandleFunc(s.saveEdit))
	mux.Handle("DELETE /delete/{id}", makeHTTPHandleFunc(s.delete))
	mux.Handle("GET /create", makeHTTPHandleFunc(s.create))
	mux.Handle("POST /create", makeHTTPHandleFunc(s.saveCreate))
	mux.Handle("GET /", makeHTTPHandleFunc(s.homepage))
	//mux.Handle("GET /cdn/{filename}", makeHTTPHandleFunc(s.debugHandler))

	err := http.ListenAndServe(s.listenAddr, mux)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *APIServer) homepage(w http.ResponseWriter, r *http.Request) {
	s.views.render(HOMEPAGE, "index", http.StatusOK, w, ALIAS)
}

func (s *APIServer) table(w http.ResponseWriter, r *http.Request) {
	sublicenses, err := s.db.GetAllSublicenses()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	s.views.render(HOMEPAGE, "table", http.StatusOK, w, sublicenses)
}

func (s *APIServer) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		InvaildRequest(w, r, "Error with ID")
	}
	err = s.db.DeleteSublicenseById(id)
	if err != nil {
		InvaildRequest(w, r, fmt.Sprintf("Error Deleting: %s", err))
	}
}

func (s *APIServer) edit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	license, err := s.db.GetSublicense(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	scrnData := SublicenseScreen{
		Alias: ALIAS,
		Data:  *license,
	}
	s.views.render(SUBLICENSE, "sublicense-edit", http.StatusOK, w, scrnData)
}

func (s *APIServer) create(w http.ResponseWriter, r *http.Request) {
	scrnData := SublicenseScreen{
		Alias: ALIAS,
		Data:  Sublicense{},
	}
	s.views.render(SUBLICENSE, "sublicense-create", http.StatusOK, w, scrnData)
}

func (s *APIServer) saveEdit(w http.ResponseWriter, r *http.Request) {

}

func (s *APIServer) saveCreate(w http.ResponseWriter, r *http.Request) {
	fname := r.FormValue("fname")
	if len(fname) > 255 || len(fname) == 0 {
		InvaildRequest(w, r, "Name to long or empty")
		return
	}
	seats := r.FormValue("fnumberOfSeats")
	if len(seats) == 0 {
		InvaildRequest(w, r, "Number of Seats empty")
		return
	}
	numberofseats, err := strconv.Atoi(seats)
	if err != nil && numberofseats > 0 {
		InvaildRequest(w, r, "Number of Seats")
		return
	}
	key := LicenseKey(r.FormValue("fLicenseKey"))
	if !key.valid() {
		InvaildRequest(w, r, "License Key invalid")
		return
	}
	date, err := HTMLDateStringToUnixtime(r.FormValue("fExpiryDate"))
	if err != nil {
		InvaildRequest(w, r, "Expiry Date")
		return
	}
	expiryDate := UnixtimeToHTMLDateString(date)
	activ := false
	if strings.ToUpper(r.FormValue("fActiv")) == "ON" {
		activ = true
	}
	lic := Sublicense{
		Id:            -1, // So I know that I have to ignore it when writing
		Name:          fname,
		NumberOfSeats: numberofseats,
		LicenseKey:    string(key),
		ExpiryDate:    expiryDate,
		Activ:         activ,
	}
	if err = s.db.CreateSublicense(&lic); err != nil {
		log.Println(err)
		InvaildRequest(w, r, "Saving Failed")
		return
	}
}

func (s *APIServer) debugHandler(_ http.ResponseWriter, r *http.Request) {
	reqDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("REQUEST:\n%s", string(reqDump))
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r)
	}
}

func loadTemplates() Views {
	var out []*template.Template
	var fileList []string = []string{
		"views/index.html",
		"views/sublicense.html",
	}
	for _, fileName := range fileList {
		tmpl, err := template.ParseFiles(fileName)
		if err != nil {
			panic(err)
		}
		out = append(out, tmpl)
	}
	return Views{
		pages: out,
	}
}

func InvaildRequest(w http.ResponseWriter, r *http.Request, extraInfo string) {
	http.Error(w, fmt.Sprintf("Error: %s", extraInfo), http.StatusBadRequest)
}
