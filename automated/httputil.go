package main

import (
	"encoding/base64"
	"net/http"
	"strings"
)

type BasicAuth struct {
	User, Password string
}

func getBasicAuth(r *http.Request) *BasicAuth {
	var userpass string
	auth := r.Header.Get("Authorization")
	if auth[:6] == "Basic " {
		userpass = auth[6:]
	}
	decoded, err := base64.StdEncoding.DecodeString(userpass)
	if err != nil {
		return nil
	}
	split := strings.Split(string(decoded), ":")
	user := split[0]
	pass := ""
	if len(split) > 1 {
		pass = split[1]
	}

	return &BasicAuth{user, pass}
}
