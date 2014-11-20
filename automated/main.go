package main

import (
	"flag"
	"log"
	"net/url"
	"regexp"
)

func main() {
	port := flag.String("port", "8888", "Plain http listen port")
	tlsport := flag.String("tlsport", "8444", "Plain http listen port")
	tls := flag.Bool("tls", false, "Enable TLS port")
	cert := flag.String("cert", "cert.pem", "TLS Certificate")
	key := flag.String("key", "key.pem", "TLS Key")
	ohurl := flag.String("ohurl", "http://localhost:8080/rest", "OpenHAB URL")
	static := flag.String("static", "/usr/local/share/automated", "Static files")
	dynamic := flag.String("var", "/var/local/automated", "Static files")
	flag.Parse()

	url, err := url.Parse(*ohurl)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("OH_URL: ", url.String())
	oldbase := url.Scheme + "://" + url.Host

	srv := &server{
		Static:  *static,
		Dynamic: *dynamic,
		Port:    *port,
		OHURL:   url.String()}

	srv.rew = regexp.MustCompile(oldbase)
	if *tls {
		srv.TLSPort = *tlsport
		srv.Key = *key
		srv.Cert = *cert
	}

	log.Fatal(srv.Run())
}
