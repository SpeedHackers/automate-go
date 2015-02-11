package main

import (
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/Pursuit92/httputil"
	"github.com/Pursuit92/syncmap"
	"github.com/gorilla/mux"
)

type server struct {
	OHURL     string
	rew       *regexp.Regexp
	Port      string
	TLSPort   string
	Cert      string
	Key       string
	Static    string
	Dynamic   string
	db        *DB
	PermCache syncmap.Map
}

func (s *server) setupRoutes() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/yo/items/{item}", s.yo).Methods("GET")
	r.HandleFunc("/rest", s.rest).
		Methods("GET")
	r.PathPrefix("/images/").
		Handler(http.StripPrefix("/images/",
		http.FileServer(http.Dir(s.Static+"/images/"))))
	r.HandleFunc("/hooks", s.hooks).Methods("POST")
	rest := r.PathPrefix("/rest").Subrouter()
	rest.HandleFunc("/", s.rest).
		Methods("GET")
	rest.HandleFunc("/items", s.getItems).
		Methods("GET")
	rest.HandleFunc("/sitemaps", s.getMaps).
		Methods("GET")

	maps := rest.PathPrefix("/sitemaps").Subrouter()
	maps.HandleFunc("/", s.getMaps).
		Methods("GET")
	maps.HandleFunc("/{map}", s.getMap).
		Methods("GET")
	maps.HandleFunc("/{map}/{page}", s.getPage).
		Methods("GET")
	maps.HandleFunc("/{map}/{page}", s.getPageStreaming).
		Headers("X-Atmosphere-Transport", "streaming").
		Methods("GET")

	items := rest.PathPrefix("/items").Subrouter()
	items.HandleFunc("/", s.getItems).
		Methods("GET")
	items.HandleFunc("/{item}", s.getItemStreaming).
		Headers("X-Atmosphere-Transport", "streaming").
		Methods("GET")
	items.HandleFunc("/{item}", s.getItem).
		Methods("GET")
	items.HandleFunc("/{item}", s.cmdItem).
		Methods("POST")

	return httputil.Logger(r)
}

func (s *server) Run() error {
	ch := make(chan error)

	ohurl, err := url.Parse(s.OHURL)
	if err != nil {
		return err
	}
	s.OHURL = ohurl.String()
	log.Print("OH_URL: ", ohurl.String())
	oldbase := ohurl.Scheme + "://" + ohurl.Host
	s.rew = regexp.MustCompile(oldbase)

	s.PermCache = syncmap.New()

	routes := s.setupRoutes()
	/*
		var err error
		s.db, err = OpenDB(s.Dynamic + "/db")
		if err != nil {
			return err
		}
	*/

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
