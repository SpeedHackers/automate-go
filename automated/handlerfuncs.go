package main

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *server) rest(w http.ResponseWriter, r *http.Request) {
	client := makeClient(r, s.OHURL)

	data, err := client.Base()
	if err != nil {
		Error(w, err)
		return
	}

	s.finish(r, w, data)
}

func (s *server) getMaps(w http.ResponseWriter, r *http.Request) {
	client := makeClient(r, s.OHURL)

	data, err := client.Sitemaps()
	if err != nil {
		Error(w, err)
		return
	}

	s.finish(r, w, data)
}

func (s *server) getMap(w http.ResponseWriter, r *http.Request) {
	client := makeClient(r, s.OHURL)

	smap := mux.Vars(r)["map"]
	data, err := client.Sitemap(smap)
	if err != nil {
		Error(w, err)
		return
	}

	s.finish(r, w, data)
}

func (s *server) getPage(w http.ResponseWriter, r *http.Request) {
	client := makeClient(r, s.OHURL)

	vars := mux.Vars(r)
	smap := vars["map"]
	page := vars["page"]
	transport := r.Header.Get("X-Atmosphere-Transport")

	switch transport {
	case "long-polling":
		pageerr := <-client.PageLongPolling(smap, page)
		if pageerr.Error != nil {
			Error(w, pageerr.Error)
			return
		}
		s.finish(r, w, pageerr.Page)
	default:
		data, err := client.SitemapPage(smap, page)
		if err != nil {
			Error(w, err)
			return
		}

		s.finish(r, w, data)
	}
}

func (s *server) getPageStreaming(w http.ResponseWriter, r *http.Request) {
	client := makeClient(r, s.OHURL)

	vars := mux.Vars(r)
	smap := vars["map"]
	page := vars["page"]

	ch, ctl := client.PageStreaming(smap, page)
	defer close(ctl)
	for pageerr := range ch {
		if pageerr.Error != nil {
			Error(w, pageerr.Error)
			return
		} else {
			data := pageerr.Page
			err := s.finish(r, w, data)
			if err != nil {
				return
			}
		}
	}
}

func (s *server) getItem(w http.ResponseWriter, r *http.Request) {
	client := makeClient(r, s.OHURL)

	item := mux.Vars(r)["item"]
	transport := r.Header.Get("X-Atmosphere-Transport")

	switch transport {
	case "long-polling":
		iterr := <-client.ItemLongPolling(item)
		if iterr.Error != nil {
			Error(w, iterr.Error)
			return
		}
		s.finish(r, w, iterr.Item)
	default:
		data, err := client.Item(item)
		if err != nil {
			Error(w, err)
			return
		}

		s.finish(r, w, data)
	}
}
func (s *server) getItemStreaming(w http.ResponseWriter, r *http.Request) {
	client := makeClient(r, s.OHURL)

	item := mux.Vars(r)["item"]
	ch, ctl := client.ItemStreaming(item)
	defer close(ctl)
	for iterr := range ch {
		if iterr.Error != nil {
			Error(w, iterr.Error)
			return
		} else {
			data := iterr.Item
			err := s.finish(r, w, data)
			if err != nil {
				return
			}
		}
	}
}

func (s *server) cmdItem(w http.ResponseWriter, r *http.Request) {
	client := makeClient(r, s.OHURL)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error(w, err)
		return
	}

	item := mux.Vars(r)["item"]
	err = client.CommandItem(item, string(body))
	if err != nil {
		Error(w, err)
		return
	}

	w.WriteHeader(201)

	s.finish(r, w, nil)
}
func (s *server) getItems(w http.ResponseWriter, r *http.Request) {
	client := makeClient(r, s.OHURL)

	data, err := client.Items()
	if err != nil {
		Error(w, err)
		return
	}

	s.finish(r, w, data)
}
