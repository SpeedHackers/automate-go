package main

/*
func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(loggerFunc(h.ServeHTTP))
}

func loggerFunc(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		reqLog := fmt.Sprintf("%s: %s %s %s", r.RemoteAddr, r.Method, r.URL.String(), r.Proto)
		resp := httputil.NewSnifferWriter(w)
		f(resp, r)
		if resp.Status == 0 {
			resp.Status = 200
		}
		log.Printf("%s %d", reqLog, resp.Status)
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
		exists, err := s.db.Exists("users", auth.Username)
		if err != nil {
			restErr := openhab.NewRestError(err)
			http.Error(w, restErr.Text, restErr.Status)
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
				http.Error(w, restErr.Text, restErr.Status)
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
