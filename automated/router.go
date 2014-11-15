package main

import (
	"net/http"

	"github.com/SpeedHackers/automate-go/openhab"
	"github.com/gorilla/mux"
)

type server struct {
	Client *openhab.Client
}

func (s *server) setupRoutes() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/rest", s.base).Methods("GET")
	rest := r.PathPrefix("/rest").Subrouter()
	rest.HandleFunc("/", s.base).Methods("GET")
	rest.HandleFunc("/items", s.getItems).Methods("GET")
	rest.HandleFunc("/sitemaps", s.getMaps).Methods("GET")

	maps := rest.PathPrefix("/sitemaps").Subrouter()
	maps.HandleFunc("/", s.getMaps).Methods("GET")
	maps.HandleFunc("/{map}", s.getMap).Methods("GET")

	items := rest.PathPrefix("/items").Subrouter()
	items.HandleFunc("/", s.getItems).Methods("GET")
	items.HandleFunc("/{item}", s.getItem).Methods("GET")
	items.HandleFunc("/{item}", s.cmdItem).Methods("POST")

	return logger(r)
}