package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

type snifferWriter struct {
	buf    bytes.Buffer
	code   int
	header http.Header
}

func newSnifferWriter() *snifferWriter {
	return &snifferWriter{header: make(map[string][]string)}
}
func (l *snifferWriter) Header() http.Header {
	return l.header
}

func (l *snifferWriter) WriteHeader(i int) {
	l.code = i
}

func (l *snifferWriter) Write(b []byte) (int, error) {
	return l.buf.Write(b)
}

func (l *snifferWriter) WriteOut(w http.ResponseWriter) {
	h := l.Header()
	nh := w.Header()
	for i, v := range h {
		nh[i] = v
	}
	if l.code == 0 {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(l.code)
	}
	_, err := io.Copy(w, &l.buf)
	if err != nil {
		log.Print("Error: ", err)
	}
}

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(loggerFunc(h.ServeHTTP))
}

func loggerFunc(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		reqLog := fmt.Sprintf("%s: %s %s %s", r.RemoteAddr, r.Method, r.URL.String(), r.Proto)
		resp := newSnifferWriter()
		f(resp, r)
		if resp.code == 0 {
			resp.code = 200
		}
		log.Printf("%s %d", reqLog, resp.code)
		resp.WriteOut(w)
	}
}

/*
func (s *server) auth(h http.Handler) http.Handler {
	return http.HandlerFunc(s.authFunc(h.ServeHTTP))
}

func (s *server) authFunc(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := newSnifferWriter()
		f(resp, r)
		auth := getBasicAuth(r)
		exists, err := s.db.Exists("users", auth.Username)
		if err != nil {
			restErr := openhab.NewRestError(err)
			http.Error(w, restErr.Text, restErr.Code)
			return
		}
		if !exists {
			s.db.Set("users", auth.Username, User{auth.Username, auth.Password, []string{auth.Username}})
		}
		user := &User{}
		s.db.Get("users", auth.Username, user)
		if user.Password != auth.Password {
			http.Error(w, "Invalid Login", 403)
			return
		}

		if r.Method == "POST" {
			parts := strings.Split(r.URL.Path, "/")
			item := parts[len(parts)-1]
			dbItem := &DBItem{}
			err := s.db.Get("items", item, dbItem)
			if err != nil {
				restErr := openhab.NewRestError(err)
				http.Error(w, restErr.Text, restErr.Code)
				return
			}

		}
		resp.WriteOut(w)
	}
}

/*
func (s *server) requireAuthFunc(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sess := s.getSession(r, w)
		sess.Last = r.URL.Path
		if sess.Authenticated {
			f(w, r)
		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	}
}
*/
