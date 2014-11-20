package openhab

// Get the base of the rest endpoint
func (cl *Client) Base() (RestBase, error) {
	resp := RestBase{}
	outch, _ := cl.request("GET", "", "", &resp, NormalReq)
	out := <-outch
	if out.Error != nil {
		return RestBase{}, out.Error
	}

	return out.Val.(RestBase), nil
}

// Get a list of Sitemaps
func (cl *Client) Sitemaps() ([]Sitemap, error) {
	resp := SitemapsResp{}
	outch, _ := cl.request("GET", "/sitemaps", "", &resp, NormalReq)
	out := <-outch
	if out.Error != nil {
		return nil, out.Error
	}

	return out.Val.(SitemapsResp).Sitemaps, nil
}

// Get a single Sitemap
func (cl *Client) Sitemap(name string) (Sitemap, error) {
	resp := Sitemap{}
	outch, _ := cl.request("GET", "/sitemaps/"+name, "", &resp, NormalReq)
	out := <-outch
	if out.Error != nil {
		return Sitemap{}, out.Error
	}

	return out.Val.(Sitemap), nil
}

// Get a sitemap page
func (cl *Client) SitemapPage(name, page string) (SitemapPage, error) {
	resp := SitemapPage{}
	outch, _ := cl.request("GET", "/sitemaps/"+name+"/"+page, "", &resp, NormalReq)
	out := <-outch
	if out.Error != nil {
		return SitemapPage{}, out.Error
	}

	return out.Val.(SitemapPage), nil
}

// Get all of the items
func (cl *Client) Items() ([]Item, error) {
	resp := ItemsResp{}
	outch, _ := cl.request("GET", "/items", "", &resp, NormalReq)
	out := <-outch
	if out.Error != nil {
		return nil, out.Error
	}

	return out.Val.(ItemsResp).Items, nil
}

// Get a single Item
func (cl *Client) Item(name string) (Item, error) {
	resp := Item{}
	outch, _ := cl.request("GET", "/items/"+name, "", &resp, NormalReq)
	out := <-outch
	if out.Error != nil {
		return Item{}, out.Error
	}

	return out.Val.(Item), nil
}

// Send a command to an item
func (cl *Client) CommandItem(item, cmd string) error {
	resp, _ := cl.request("POST", "/items/"+item, cmd, nil, NormalReq)
	return (<-resp).Error
}

// Update the state of an item. Not really sure what this is for.
func (cl *Client) UpdateItem(item, cmd string) error {
	resp, _ := cl.request("PUT", "/items/"+item+"/state", cmd, nil, NormalReq)
	return (<-resp).Error
}
