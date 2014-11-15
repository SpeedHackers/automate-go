package main

import (
	"bytes"
	"net/http"
	"io"
	"log"
	"fmt"
)

type loggerWriter struct {
	buf bytes.Buffer
	code int
	header http.Header
}

func newLoggerWriter() *loggerWriter {
	return &loggerWriter{header: make(map[string][]string)}
}
func (l *loggerWriter) Header() http.Header {
	return l.header
}

func (l *loggerWriter) WriteHeader(i int) {
	l.code = i
}

func (l *loggerWriter) Write(b []byte) (int,error) {
	return l.buf.Write(b)
}

func (l *loggerWriter) WriteOut(w http.ResponseWriter) {
	h := l.Header()
	nh := w.Header()
	for i,v := range h {
		nh[i] = v
	}
	if l.code == 0 {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(l.code)
	}
	_,err := io.Copy(w,&l.buf)
	if err != nil {
		log.Print("Error: ",err)
	}
}

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(loggerFunc(h.ServeHTTP))
}

func loggerFunc(f func(http.ResponseWriter,*http.Request)) func(http.ResponseWriter,*http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		reqLog := fmt.Sprintf("%s: %s %s %s",r.RemoteAddr, r.Method,r.URL.String(),r.Proto)
		resp := newLoggerWriter()
		f(resp,r)
		if resp.code == 0 {
			resp.code = 200
		}
		log.Printf("%s %d",reqLog,resp.code )
		resp.WriteOut(w)
	}
}

