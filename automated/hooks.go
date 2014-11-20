package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	client := makeClient(r, s.OHURL)
	cmd := &Hook{}
	err := json.NewDecoder(r.Body).Decode(cmd)
	if err != nil {
		Error(w, err)
		return
	}
	log.Print("received hook: ", cmd)
	for _, v := range cmd.Cmds {
		err = client.CommandItem(v.Item, v.Cmd)
		if err != nil {
			Error(w, err)
		}
	}
	w.WriteHeader(201)
}

func (s *server) yo(w http.ResponseWriter, r *http.Request) {
	client := makeClient(r, s.OHURL)
	name := mux.Vars(r)["item"]
	item, err := client.Item(name)
	if err != nil {
		Error(w, err)
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
		err := client.CommandItem(name, cmd)
		if err != nil {
			Error(w, err)
			return
		}
	case "StringItem":
		loc, ok := r.Form["location"]
		if ok {
			client.CommandItem(name, loc[0])
		} else {
			user := r.Form["username"]
			client.CommandItem(name, user[0])
		}

	case "NumberItem":
		n, err := strconv.Atoi(item.State)
		if err != nil {
			n = 0
		} else {
			n++
		}
		err = client.CommandItem(name, fmt.Sprintf("%d", n))
		if err != nil {
			Error(w, err)
			return
		}
	default:
		http.Error(w, "Can't yo that item", 405)
		return
	}
	w.WriteHeader(201)
}
