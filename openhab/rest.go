package openhab

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type Client struct {
	Username, Password string
	URL                string
	httpClient         *http.Client
}

// Creates a new OpenHAB client. Expects a path to the start of the REST
// endpoint, e.g. "http://example.com:8000/rest"
// The sslverify argument optionally disables ssl checking (probably required)
func NewClient(url, user, pass string, sslverify bool) *Client {
	httpClient := &http.Client{}
	if !sslverify {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	cl := &Client{Username: user,
		Password:   pass,
		URL:        url,
		httpClient: httpClient}

	return cl
}

type ReqType byte

const (
	NormalReq ReqType = iota
	LongPolling
	Streaming
)

// Needs to be updated for long-polling/streaming reqs (add a chan?)
func (cl *Client) request(method, url, body string, out interface{}, reqType ReqType) (err error) {
	var req *http.Request
	if body != "" {
		bodyBuffer := bytes.NewBuffer([]byte(body))
		req, err = http.NewRequest(method, cl.URL+url, bodyBuffer)
		if err != nil {
			return
		}
	} else {
		req, err = http.NewRequest(method, cl.URL+url, nil)
		if err != nil {
			return
		}
	}
	req.Header.Add("Content-Type", "text/plain")
	req.Header.Add("Accept", "application/json")
	if cl.Username != "" && cl.Password != "" {
		req.SetBasicAuth(cl.Username, cl.Password)
	}
	resp, err := cl.httpClient.Do(req)
	if err != nil {
		return
	}
	if resp.Status[0] == '2' {
		if out != nil {
			decoder := json.NewDecoder(resp.Body)
			err = decoder.Decode(out)
		}
		return
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = errors.New(resp.Status + ": " + string(bodyBytes))
	return
}

// Get a list of Sitemaps
func (cl *Client) Sitemaps() ([]Sitemap, error) {
	resp := SitemapsResp{}
	err := cl.request("GET", "/sitemaps", "", &resp, NormalReq)
	if err != nil {
		return nil, err
	}

	return resp.Sitemaps, nil
}

// Get a single Sitemap
func (cl *Client) Sitemap(name string) (*Sitemap, error) {
	resp := Sitemap{}
	err := cl.request("GET", "/sitemaps/"+name, "", &resp, NormalReq)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Get a sitemap page
func (cl *Client) SitemapPage(name, page string) (*SitemapPage, error) {
	resp := SitemapPage{}
	err := cl.request("GET", "/sitemaps/"+name+"/"+page, "", &resp, NormalReq)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Get all of the items
func (cl *Client) Items() ([]Item, error) {
	resp := ItemsResp{}
	err := cl.request("GET", "/items", "", &resp, NormalReq)
	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}

// Get a single Item
func (cl *Client) Item(name string) (*Item, error) {
	resp := Item{}
	err := cl.request("GET", "/items/"+name, "", &resp, NormalReq)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Send a command to an item
func (cl *Client) CommandItem(item, cmd string) error {
	return cl.request("POST", "/items/"+item, cmd, nil, NormalReq)
}

// Update the state of an item. Not really sure what this is for.
func (cl *Client) UpdateItem(item, cmd string) error {
	return cl.request("PUT", "/items/"+item, cmd, nil, NormalReq)
}

// Stub for long-polling item
func (cl *Client) ItemLongPolling(item string) (chan Item, error) {
	return nil, nil
}

// Stub for streaming item
func (cl *Client) ItemStreaming(item string) (chan Item, error) {
	return nil, nil
}
