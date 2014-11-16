package main

import (
	"net/http"
	"regexp"

	"github.com/SpeedHackers/automate-go/openhab"
	"github.com/gorilla/mux"
)

type server struct {
	Client  *openhab.Client
	rew     *regexp.Regexp
	Port    string
	TLSPort string
	db      *DB
}

func (s *server) setupRoutes() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/yo/items/{item}", s.yo).Methods("GET")
	r.HandleFunc("/rest", s.base).Methods("GET")
	r.HandleFunc("/hooks", s.hooks).Methods("POST")
	r.HandleFunc("/db", s.dbHandler).Methods("POST", "GET")
	rest := r.PathPrefix("/rest").Subrouter()
	rest.HandleFunc("/", s.base).Methods("GET")
	rest.HandleFunc("/items", s.getItems).Methods("GET")
	rest.HandleFunc("/sitemaps", s.getMaps).Methods("GET")

	maps := rest.PathPrefix("/sitemaps").Subrouter()
	maps.HandleFunc("/", s.getMaps).Methods("GET")
	maps.HandleFunc("/{map}", s.getMap).Methods("GET")
	maps.HandleFunc("/{map}/{page}", s.getPage).Methods("GET")

	items := rest.PathPrefix("/items").Subrouter()
	items.HandleFunc("/", s.getItems).Methods("GET")
	items.HandleFunc("/{item}", s.getItem).Methods("GET")
	items.HandleFunc("/{item}", s.cmdItem).Methods("POST")

	return logger(s.rewriter(r))
}

func (s *server) Run() error {
	ch := make(chan error)
	routes := s.setupRoutes()
	var err error
	s.db, err = OpenDB(".")
	if err != nil {
		return err
	}

	go func() {
		ch <- http.ListenAndServe(":"+s.Port, routes)
	}()
	return <-ch
}
