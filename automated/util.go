package main

import (
	"strings"

	"github.com/SpeedHackers/automate-go/openhab"
)

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

			}
		}
	}
	topGrp.Members = nil
	items = append(items, topGrp)
	return items, nil
}

func getAllowed(cl *openhab.Client) ([]openhab.Item, error) {
	allowed, err := getGroupRec(cl, "Group_"+strings.Title(cl.Username))
	return allowed, err
}
