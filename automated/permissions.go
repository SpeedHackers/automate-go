package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/SpeedHackers/automate-go/openhab"
)

type User struct {
	Username, Password string
	FakeUser, FakePass string
	Items              openhab.Items
	Expire             time.Time
}

func inItems(it string, its []openhab.Item) bool {
	for _, v := range its {
		if v.Name == it {
			return true
		}
	}
	return false
}

func getGroupRec(cl ohClient, name string) (openhab.Items, error) {
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

func (s *server) getAllowed(cl ohClient) ([]openhab.Item, error) {
	var usrInt interface{}
	var ok bool
	if cl.FakeUser != "" {
		usrInt, ok = s.PermCache.Get(cl.FakeUser)
		if ok {
			usr := usrInt.(User)
			return usr.Items, nil
		}
		return nil, fmt.Errorf("User not in cache")
	}
	usrInt, ok = s.PermCache.Get(cl.Username)
	if ok {
		usr := usrInt.(User)
		if usr.Password != cl.Password {
			return nil, fmt.Errorf("Credentials differ from cache")
		}
		return usr.Items, nil
	}
	allowed, err := getGroupRec(cl, "Group_"+strings.Title(cl.Username))
	if err == nil {
		usr := User{Username: cl.Username,
			Password: cl.Password,
			Items:    allowed,
			Expire:   time.Now().Add(1 * time.Hour)}
		s.PermCache.Set(cl.Username, usr)
		s.startRefresh(cl, usr)
	}
	return allowed, err
}

func (s *server) startRefresh(cl ohClient, usr User) {
	go func() {
		for {
			<-time.After(1 * time.Minute)
			if time.Now().After(usr.Expire) {
				s.PermCache.Delete(cl.Username)
				return
			}
			allowed, err := getGroupRec(cl, "Group_"+strings.Title(cl.Username))
			if err == nil {
				usr.Items = allowed
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
