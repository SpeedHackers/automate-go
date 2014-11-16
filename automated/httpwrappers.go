package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/SpeedHackers/automate-go/openhab"
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
		log.Printf("%s %d", reqLog, resp.code)
		//log.Print(r.Header)
		resp.WriteOut(w)
	}
}

func (s *server) rewriter(h http.Handler) http.Handler {
	return http.HandlerFunc(s.rewriteFunc(h.ServeHTTP))
}

func (s *server) rewriteFunc(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := newSnifferWriter()
		f(resp, r)
		var scheme string
		hostport := strings.Split(r.Host, ":")
		if len(hostport) > 1 && hostport[1] == s.TLSPort {
			scheme = "https"
		} else {
			scheme = "http"
		}
		body := resp.buf.Bytes()
		resp.buf = *bytes.NewBuffer(s.rew.ReplaceAll(body, []byte(scheme+"://"+r.Host)))
		resp.header.Add("Content-Type", "application/json")
		resp.header.Add("Access-Control-Allow-Origin", "*")
		resp.WriteOut(w)
	}
}

func (s *server) auth(h http.Handler) http.Handler {
	return http.HandlerFunc(s.authFunc(h.ServeHTTP))
}

func (s *server) authFunc(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := newSnifferWriter()
		f(resp, r)
		auth := getBasicAuth(r)
		exists, err := s.db.Exists("users", auth.User)
		if err != nil {
			restErr := openhab.NewRestError(err)
			http.Error(w, restErr.Text, restErr.Code)
			return
		}
		if !exists {
			s.db.Set("users", auth.User, User{auth.User, auth.Password, []string{auth.User}})
		}
		user := &User{}
		s.db.Get("users", auth.User, user)
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
