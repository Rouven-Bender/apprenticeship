package main

import (
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
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
		Conversions: SublicenseConverions{
			Date: UnixtimeToHTMLDateString(license.ExpiryDate),
		},
	}
	s.views.render(SUBLICENSE, "sublicense-edit", http.StatusOK, w, scrnData)
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
