package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/SpeedHackers/automate-go/openhab"
	"github.com/gorilla/mux"
)

func (s *server) base(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "items\nsitemaps\n")
}

func (s *server) getMaps(w http.ResponseWriter, r *http.Request) {
	maps, err := s.Client.Sitemaps()
	if err != nil {
		oherr := err.(openhab.RestError)
		http.Error(w, oherr.Text, oherr.Code)
		return
	}
	err = json.NewEncoder(w).Encode(maps)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
func (s *server) getMap(w http.ResponseWriter, r *http.Request) {
	smap, err := s.Client.Sitemap(mux.Vars(r)["map"])
	if err != nil {
		oherr := err.(openhab.RestError)
		http.Error(w, oherr.Text, oherr.Code)
		return
	}
	err = json.NewEncoder(w).Encode(smap)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
func (s *server) getItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	item, err := s.Client.Item(vars["item"])
	if err != nil {
		oherr := err.(openhab.RestError)
		http.Error(w, oherr.Text, oherr.Code)
		return
	}
	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
func (s *server) cmdItem(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = s.Client.CommandItem(mux.Vars(r)["item"], string(body))
	if err != nil {
		oherr := err.(openhab.RestError)
		http.Error(w, oherr.Text, oherr.Code)
		return
	}
}
func (s *server) getItems(w http.ResponseWriter, r *http.Request) {
	items, err := s.Client.Items()
	if err != nil {
		oherr := err.(openhab.RestError)
		http.Error(w, oherr.Text, oherr.Code)
		return
	}
	err = json.NewEncoder(w).Encode(items)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
