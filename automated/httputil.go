package main

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/SpeedHackers/automate-go/openhab"
)

type BasicAuth struct {
	Username, Password string
}

type ohClient struct {
	*openhab.Client
	FakeUser string
}

func getBasicAuth(r *http.Request) BasicAuth {
	var userpass string
	auth := r.Header.Get("Authorization")
	if len(auth) == 0 {
		return BasicAuth{}
	}
	if auth[:6] == "Basic " {
		userpass = auth[6:]
	}
	decoded, err := base64.StdEncoding.DecodeString(userpass)
	if err != nil {
		return BasicAuth{}
	}
	split := strings.Split(string(decoded), ":")
	user := split[0]
	pass := ""
	if len(split) > 1 {
		pass = split[1]
	}

	return BasicAuth{user, pass}
}

func getEncoding(r *http.Request) (Encoding, string) {
	encHeader := r.Header["Accept"]
	var encode Encoding
	var str string
	if len(encHeader) > 0 {
		switch encHeader[0] {
		case "", "application/json", "text/json":
			str = "json"
			encode = JSONEncoding
		default:
			str = "xml"
			encode = XMLEncoding
		}
	} else {
		str = "json"
		encode = JSONEncoding
	}
	return encode, str
}

func (s *server) makeClient(r *http.Request, url string) ohClient {
	auth := getBasicAuth(r)
	username := auth.Username
	password := auth.Password
	fakeUser := ""
	usrInt, ok := s.PermCache.Get(username)
	if ok {
		usr := usrInt.(User)
		if usr.FakeUser == username && usr.FakePass == password {
			fakeUser = username
			username = usr.Username
			password = usr.Password
		}
	}

	return ohClient{openhab.NewClient(url, username, password, false), fakeUser}
}

func (s *server) finish(r *http.Request, w http.ResponseWriter, data interface{}) error {
	encode, str := getEncoding(r)

	// Set correct headers
	head := w.Header()
	if _, ok := head["Content-Type"]; !ok {
		if str == "json" {
			head.Add("Content-Type", "application/json")
		} else {
			head.Add("Content-Type", "text/xml")
		}
	}
	if _, ok := head["Access-Control-Allow-Origin"]; !ok {
		head.Add("Access-Control-Allow-Origin", "*")
	}

	// rewrite urls from openhab
	var scheme string
	hostport := strings.Split(r.Host, ":")
	if len(hostport) > 1 && hostport[1] == s.TLSPort {
		scheme = "https"
	} else {
		scheme = "http"
	}
	tmp := &bytes.Buffer{}
	err := encode(tmp, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return err
	}

	w.Write(s.rew.ReplaceAll(tmp.Bytes(), []byte(scheme+"://"+r.Host)))
	return err
}

func Error(w http.ResponseWriter, err error) {
	if ohErr, ok := err.(openhab.RestError); ok {
		if ohErr.Code == 401 {
			head := w.Header()
			head.Add("WWW-Authenticate", "Basic realm=\"openhab.org\"")
			head.Add("Content-Type", "text/html")
		}
		http.Error(w, ohErr.Text, ohErr.Code)
	} else {
		http.Error(w, err.Error(), 500)
	}
}
