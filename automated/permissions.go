package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/SpeedHackers/automate-go/openhab"
)

type User struct {
	Userame, Password string
	Items             openhab.Items
}

func inItems(it string, its []openhab.Item) bool {
	for _, v := range its {
		if v.Name == it {
			return true
		}
	}
	return false
}

func getGroupRec(cl *openhab.Client, name string) (openhab.Items, error) {
	topGrp, err := cl.Item(name)
	if err != nil {
		return nil, err
	}
	var items openhab.Items
	if topGrp.Members != nil {
		for _, v := range topGrp.Members {
			if v.Type == "GroupItem" {
				subItems, err := getGroupRec(cl, v.Name)
				if err != nil {
					return nil, err
				}
				items = append(items, subItems...)
			} else {
				items = append(items, v)
			}
		}
	}
	topGrp.Members = nil
	items = append(items, topGrp)
	return items, nil
}

func (s *server) getAllowed(cl *openhab.Client) ([]openhab.Item, error) {
	usrInt, ok := s.PermCache.Get(cl.Username)
	if ok {
		usr := usrInt.(User)
		if usr.Password != cl.Password {
			return nil, fmt.Errorf("Credentials differ from cache")
		}
		return usr.Items, nil
	}
	allowed, err := getGroupRec(cl, "Group_"+strings.Title(cl.Username))
	if err == nil {
		usr := User{cl.Username, cl.Password, allowed}
		s.PermCache.Set(cl.Username, usr)
		s.startRefresh(cl)
	}
	return allowed, err
}

func (s *server) startRefresh(cl *openhab.Client) {
	go func() {
		for {
			<-time.After(30 * time.Second)
			allowed, err := getGroupRec(cl, "Group_"+strings.Title(cl.Username))
			if err == nil {
				usr := User{cl.Username, cl.Password, allowed}
				s.PermCache.Set(cl.Username, usr)
			} else {
				s.PermCache.Delete(cl.Username)
				return
			}
		}
	}()
}

func filterPage(page *openhab.SitemapPage, allowed openhab.Items) {
	var wsNew openhab.Widgets
	if page.Widgets != nil {
		for _, v := range page.Widgets {
			if v.Item != nil {
				if inItems(v.Item.Name, allowed) {
					wsNew = append(wsNew, v)
				}
			} else if v.LinkedPage != nil {
				filterPage(v.LinkedPage, allowed)
				if v.LinkedPage.Widgets != nil {
					if len(v.LinkedPage.Widgets) != 0 {
						wsNew = append(wsNew, v)
					}
				} else {
					wsNew = append(wsNew, v)
				}
			}
		}
	}
	page.Widgets = wsNew
}
