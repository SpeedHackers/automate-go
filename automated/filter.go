package main

import "github.com/SpeedHackers/automate-go/openhab"

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
