package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"

	_ "embed"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

//go:embed JWT_SECRET
var JWT_SECRET []byte

//go:embed views/*.tmpl
var TEMPLATE_FS embed.FS

func (v *Views) render(p page, block string, statusCode int, w http.ResponseWriter, data any) error {
	w.WriteHeader(statusCode)
	v.views[p].ExecuteTemplate(w, block, data)
	return nil
}

func NewAPIServer(listenAddr string, database *sqliteStore) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		db:         database,
		views:      loadTemplates(),
	}
}

func (s *APIServer) Run() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /table", s.requiresAuthToken(s.table))
	mux.Handle("GET /edit/{id}", s.requiresAuthToken(s.edit))
	mux.Handle("PATCH /edit/{id}", s.requiresAuthToken(s.saveEdit))
	mux.Handle("DELETE /delete/{id}", s.requiresAuthToken(s.delete))
	mux.Handle("GET /create", s.requiresAuthToken(s.create))
	mux.Handle("POST /create", s.requiresAuthToken(s.saveCreate))
	mux.Handle("GET /home", s.requiresAuthToken(s.homepage))

	mux.Handle("GET /cdn/{filename}", http.StripPrefix("/cdn/", fs))
	mux.HandleFunc("GET /", s.index)
	mux.HandleFunc("GET /login", s.login)
	mux.HandleFunc("POST /login", s.verifyCredentials)
	//mux.Handle("GET /cdn/{filename}", makeHTTPHandleFunc(s.debugHandler))

	err := http.ListenAndServe(s.listenAddr, mux)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *APIServer) homepage(w http.ResponseWriter, r *http.Request) {
	if err := s.views.render(HOMEPAGE, "index", http.StatusOK, w, ALIAS); err != nil {
		log.Println(err)
	}
}

func (s *APIServer) index(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("authToken"); err != nil {
		s.login(w, r)
		return
	} else {
		s.homepage(w, r)
		return
	}
}

func (s *APIServer) verifyCredentials(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	pwd := r.FormValue("password")
	dbHash, _ := s.db.GetHashForUser(username)
	result := bcrypt.CompareHashAndPassword(dbHash, []byte(pwd))
	if result == nil {
		claims := &jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * 2)},
			Issuer:    "abschlussprojekt",
			Subject:   username,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(JWT_SECRET)
		if err != nil {
			InternalError(w, fmt.Sprintf("Signing Token: %s", err))
			return
		}
		cookie := http.Cookie{
			Name:    "authToken",
			Value:   tokenString,
			Quoted:  false,
			Expires: time.Now().Add(time.Hour * 2),
		}
		if err := cookie.Valid(); err != nil {
			log.Println(err)
		}
		http.SetCookie(w, &cookie)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}
}

func (s *APIServer) table(w http.ResponseWriter, r *http.Request) {
	sublicenses, err := s.db.GetAllSublicenses()
	if err != nil {
		InternalError(w, "Getting sublicenses failed")
	}
	if err := s.views.render(HOMEPAGE, "table", http.StatusOK, w, sublicenses); err != nil {
		log.Println(err)
	}
}

func (s *APIServer) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
	}
	err = s.db.DeleteSublicenseById(id)
	if err != nil {
		InvaildRequest(w, fmt.Sprintf("Error Deleting: %s", err))
	}
}

func (s *APIServer) login(w http.ResponseWriter, r *http.Request) {
	if err := s.views.render(LOGIN, "login", http.StatusOK, w, 0); err != nil {
		log.Println(err)
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
		InvaildRequest(w, fmt.Sprintf("Getting the sublicense: %s", err))
		return
	}
	scrnData := SublicenseScreen{
		Alias: ALIAS,
		Data:  *license,
	}
	if err := s.views.render(SUBLICENSE, "sublicense-edit", http.StatusOK, w, scrnData); err != nil {
		log.Println(err)
	}
}

func (s *APIServer) create(w http.ResponseWriter, r *http.Request) {
	scrnData := SublicenseScreen{
		Alias: ALIAS,
		Data:  Sublicense{},
	}
	if err := s.views.render(SUBLICENSE, "sublicense-create", http.StatusOK, w, scrnData); err != nil {
		log.Println(err)
	}
}

func (s *APIServer) saveEdit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fname := r.FormValue("fname")
	if len(fname) > 255 || len(fname) == 0 {
		InvaildRequest(w, "Name to long or empty")
		return
	}
	seats := r.FormValue("fnumberOfSeats")
	if len(seats) == 0 {
		InvaildRequest(w, "Number of Seats empty")
		return
	}
	numberofseats, err := strconv.Atoi(seats)
	if err != nil && numberofseats > 0 {
		InvaildRequest(w, "Number of Seats")
		return
	}
	key := LicenseKey(r.FormValue("fLicenseKey"))
	if !key.valid() {
		InvaildRequest(w, "License Key invalid")
		return
	}
	date, err := HTMLDateStringToUnixtime(r.FormValue("fExpiryDate"))
	if err != nil {
		InvaildRequest(w, "Expiry Date")
		return
	}
	expiryDate := UnixtimeToHTMLDateString(date)
	activ := false
	if strings.ToUpper(r.FormValue("fActiv")) == "ON" {
		activ = true
	}
	lic := &Sublicense{
		Id:            id,
		Name:          fname,
		NumberOfSeats: numberofseats,
		LicenseKey:    string(key),
		ExpiryDate:    expiryDate,
		Activ:         activ,
	}
	err = s.db.UpdateSublicense(lic)
	if err != nil {
		InvaildRequest(w, fmt.Sprintf("writing to DB: %s", err))
		return
	}
}

func (s *APIServer) saveCreate(w http.ResponseWriter, r *http.Request) {
	fname := r.FormValue("fname")
	if len(fname) > 255 || len(fname) == 0 {
		InvaildRequest(w, "Name to long or empty")
		return
	}
	seats := r.FormValue("fnumberOfSeats")
	if len(seats) == 0 {
		InvaildRequest(w, "Number of Seats empty")
		return
	}
	numberofseats, err := strconv.Atoi(seats)
	if err != nil && numberofseats > 0 {
		InvaildRequest(w, "Number of Seats")
		return
	}
	key := LicenseKey(r.FormValue("fLicenseKey"))
	if !key.valid() {
		InvaildRequest(w, "License Key invalid")
		return
	}
	date, err := HTMLDateStringToUnixtime(r.FormValue("fExpiryDate"))
	if err != nil {
		InvaildRequest(w, "Expiry Date")
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
		InvaildRequest(w, "Saving Failed")
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

func (s *APIServer) requiresAuthToken(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("authToken")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		token, err := validateJWT(cookie.Value)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if !token.Valid {
			log.Println(err)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		f(w, r)
	}
}

func loadTemplates() *Views {
	var out []*template.Template
	var fileList []string = []string{
		"views/index.tmpl",
		"views/sublicense.tmpl",
		"views/login.tmpl",
	}
	for _, filename := range fileList {
		tmp, err := template.ParseFS(TEMPLATE_FS, filename)
		if err != nil {
			log.Fatalln(err)
		}
		out = append(out, tmp)
	}
	return &Views{
		views: out,
	}
}

func InvaildRequest(w http.ResponseWriter, extraInfo string) {
	http.Error(w, fmt.Sprintf("Error: %s", extraInfo), http.StatusBadRequest)
}

func InternalError(w http.ResponseWriter, extraInfo string) {
	http.Error(w, fmt.Sprintf("Error: %s", extraInfo), http.StatusInternalServerError)
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() == "HS256" {
			if issuer, err := token.Claims.GetIssuer(); err == nil && issuer == "abschlussprojekt" {
				return JWT_SECRET, nil
			}
		}
		return nil, fmt.Errorf("Couldn't Parse Token")
	})
}
