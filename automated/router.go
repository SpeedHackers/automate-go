package main

import (
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

type server struct {
	OHURL   string
	rew     *regexp.Regexp
	Port    string
	TLSPort string
	Cert    string
	Key     string
	Static  string
	Dynamic string
	db      *DB
}

func (s *server) setupRoutes() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/yo/items/{item}", loggerFunc(s.yo)).Methods("GET")
	r.HandleFunc("/rest", s.rest).
		Methods("GET")
	r.PathPrefix("/images/").
		Handler(http.StripPrefix("/images/",
		http.FileServer(http.Dir(s.Static+"/images/"))))
	r.HandleFunc("/hooks", s.hooks).Methods("POST")
	rest := r.PathPrefix("/rest").Subrouter()
	rest.HandleFunc("/", loggerFunc(s.rest)).
		Methods("GET")
	rest.HandleFunc("/items", loggerFunc(s.getItems)).
		Methods("GET")
	rest.HandleFunc("/sitemaps", loggerFunc(s.getMaps)).
		Methods("GET")

	maps := rest.PathPrefix("/sitemaps").Subrouter()
	maps.HandleFunc("/", loggerFunc(s.getMaps)).
		Methods("GET")
	maps.HandleFunc("/{map}", loggerFunc(s.getMap)).
		Methods("GET")
	maps.HandleFunc("/{map}/{page}", loggerFunc(s.getPage)).
		Methods("GET")
	maps.HandleFunc("/{map}/{page}", s.getPageStreaming).
		Headers("X-Atmosphere-Transport", "streaming").
		Methods("GET")

	items := rest.PathPrefix("/items").Subrouter()
	items.HandleFunc("/", loggerFunc(s.getItems)).
		Methods("GET")
	items.HandleFunc("/{item}", s.getItemStreaming).
		Headers("X-Atmosphere-Transport", "streaming").
		Methods("GET")
	items.HandleFunc("/{item}", loggerFunc(s.getItem)).
		Methods("GET")
	items.HandleFunc("/{item}", loggerFunc(s.cmdItem)).
		Methods("POST")

	return r
}

func (s *server) Run() error {
	ch := make(chan error)
	routes := s.setupRoutes()
	var err error
	s.db, err = OpenDB(s.Dynamic + "/db")
	if err != nil {
		return err
	}

	go func() {
		ch <- http.ListenAndServe(":"+s.Port, routes)
	}()
	if s.TLSPort != "" && s.Key != "" && s.Cert != "" {
		go func() {
			ch <- http.ListenAndServeTLS(":"+s.TLSPort, s.Cert, s.Key, routes)
		}()
	}
	return <-ch
}
