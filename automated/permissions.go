package main

import "log"

type DBItem struct {
	Name  string
	Owner string
	Group string
	Perm  Permissions
}

type Permissions struct {
	User  bool
	Group bool
	Other bool
}

type User struct {
	Name, Password string
	Groups         []string
}

func (s *server) initDB() error {
	log.Print("Initializing DB")
	items, err := s.Client.Items()
	if err != nil {
		return err
	}

	count := 0
	added := 0
	for _, v := range items {
		count++
		if ok, err := s.db.Exists("items", v.Name); !ok && err == nil {
			added++
			dbEntry := DBItem{Name: v.Name,
				Owner: "root",
				Group: "root",
				Perm:  Permissions{true, true, false}}

			err := s.db.Set("items", v.Name, dbEntry)
			if err != nil {
				return err
			}
		}
	}
	log.Print("Added ", added, " new items")
	log.Print("Loaded ", count, " items")

	ok, err := s.db.Exists("users", "root")
	if err != nil {
		return err
	}
	if !ok {
		s.db.Set("users", "root", User{"root", "root", []string{"root"}})
		log.Print("Created default root user with password \"root\"")
	}

	return nil
}
