package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
func (s *server) getPage(w http.ResponseWriter, r *http.Request) {
	transport := r.Header.Get("X-Atmosphere-Transport")
	switch transport {
	case "streaming":
		ch, ctl := s.Client.PageStreaming(mux.Vars(r)["map"], mux.Vars(r)["page"])
		defer close(ctl)
		for smap := range ch {
			err := json.NewEncoder(w).Encode(smap)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}
	case "long-polling":
		ch := s.Client.PageLongPolling(mux.Vars(r)["map"], mux.Vars(r)["page"])
		smap := <-ch
		err := json.NewEncoder(w).Encode(smap)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	default:
		smap, err := s.Client.SitemapPage(mux.Vars(r)["map"], mux.Vars(r)["page"])
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
}
func (s *server) getItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transport := r.Header.Get("X-Atmosphere-Transport")
	switch transport {
	case "streaming":
		ch, ctl := s.Client.ItemStreaming(vars["item"])
		defer close(ctl)
		for item := range ch {
			err := json.NewEncoder(w).Encode(item)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}
	case "long-polling":
		ch := s.Client.ItemLongPolling(vars["item"])
		item := <-ch
		err := json.NewEncoder(w).Encode(item)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	default:
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

func (s *server) dbHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		data := make(map[string]interface{})
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), 400)
			bodyBytes, _ := ioutil.ReadAll(r.Body)
			log.Print(err, string(bodyBytes))
			return
		}
		err = s.db.Set("data", "loldata", data)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	case "GET":
		data := make(map[string]interface{})
		err := s.db.Get("data", "loldata", &data)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		dataBytes, _ := json.MarshalIndent(data, "", "  ")
		w.Header().Add("Content-Type", "application/json")
		w.Write(dataBytes)
	}
}
