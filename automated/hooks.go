package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/SpeedHackers/automate-go/openhab"
	"github.com/gorilla/mux"
)

type HookCmds []HookCmd
type HookCmd struct {
	Item string `json:"item"`
	Cmd  string `json:"command"`
}

type Hook struct {
	Cmds HookCmds `json:"openhab"`
}

func (i *HookCmds) UnmarshalJSON(bs []byte) error {
	single := HookCmd{}
	multiple := make([]HookCmd, 0)
	err := json.Unmarshal(bs, &multiple)
	if err != nil {
		err := json.Unmarshal(bs, &single)
		if err != nil {
			return err
		}
		*i = []HookCmd{single}
		return nil
	}
	*i = multiple
	return nil
}

func (s *server) hooks(w http.ResponseWriter, r *http.Request) {
	cmd := &Hook{}
	err := json.NewDecoder(r.Body).Decode(cmd)
	if err != nil {
		log.Print("Error decoding json")
		http.Error(w, err.Error(), 500)
		return
	}
	log.Print("received hook: ", cmd)
	for _, v := range cmd.Cmds {
		err = s.Client.CommandItem(v.Item, v.Cmd)
		if err != nil {
			oherr := err.(openhab.RestError)
			http.Error(w, oherr.Text, oherr.Code)
		}
	}
}
func (s *server) yo(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["item"]
	item, err := s.Client.Item(name)
	if err != nil {
		oherr := err.(openhab.RestError)
		http.Error(w, oherr.Text, oherr.Code)
		return
	}
	r.ParseForm()
	switch item.Type {
	case "SwitchItem":
		var cmd string
		if item.State == "ON" {
			cmd = "OFF"
		} else {
			cmd = "ON"
		}
		err := s.Client.CommandItem(name, cmd)
		if err != nil {
			oherr := err.(openhab.RestError)
			http.Error(w, oherr.Text, oherr.Code)
			return
		}
	case "StringItem":
		loc, ok := r.Form["location"]
		if ok {
			s.Client.CommandItem(name, loc[0])
		} else {
			user := r.Form["username"]
			s.Client.CommandItem(name, user[0])
		}

	case "NumberItem":
		n, err := strconv.Atoi(item.State)
		if err != nil {
			n = 0
		} else {
			n++
		}
		err = s.Client.CommandItem(name, fmt.Sprintf("%d", n))
		if err != nil {
			oherr := err.(openhab.RestError)
			http.Error(w, oherr.Text, oherr.Code)
			return
		}
	default:
		http.Error(w, "Can't yo that item", 405)
		return

	}
}
