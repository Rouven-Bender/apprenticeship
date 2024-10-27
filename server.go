package main

import (
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
)

type APIServer struct {
	listenAddr string
	views      Views
}

type Views struct {
	pages []*template.Template
}

type page int

const (
	HOMEPAGE page = 0
)

func (v *Views) render(p page, name string, status int, w http.ResponseWriter, d interface{}) error {
	w.WriteHeader(status)
	return v.pages[p].ExecuteTemplate(w, name, d)
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		views:      loadTemplates(),
	}
}

func (s *APIServer) Run() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /cdn/{filename}", http.StripPrefix("/cdn/", fs))
	mux.Handle("GET /", makeHTTPHandleFunc(s.homepage))
	//mux.Handle("GET /cdn/{filename}", makeHTTPHandleFunc(s.debugHandler))

	err := http.ListenAndServe(s.listenAddr, mux)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *APIServer) homepage(w http.ResponseWriter, r *http.Request) {
	s.views.render(HOMEPAGE, "index", http.StatusOK, w, 0)
}

func (s *APIServer) debugHandler(_ http.ResponseWriter, r *http.Request) {
	reqDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("REQUEST:\n%s", string(reqDump))
}

type apiFunc func(http.ResponseWriter, *http.Request)

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r)
	}
}

func loadTemplates() Views {
	var out []*template.Template
	tmpl, err := template.ParseFiles("views/index.html")
	if err != nil {
		panic(err)
	}
	out = append(out, tmpl)
	return Views{
		pages: out,
	}
}
