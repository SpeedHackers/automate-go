package openhab

import "fmt"

type ItemError struct {
	Item  Item
	Error error
}
type PageError struct {
	Page  SitemapPage
	Error error
}

func (cl *Client) streamPath(path string, t interface{}) (chan interface{}, chan struct{}) {
	ch := make(chan interface{})
	ctl := make(chan struct{})
	go func() {
		outch, rctl := cl.request("GET", path, "", t, Streaming)
		for {
			select {
			case <-ctl:
				close(ch)
				close(rctl)
				return
			case out := <-outch:

				if out.Error != nil {
					close(ch)
					close(rctl)
					return
				}
				select {
				case <-ctl:
					close(ch)
					close(rctl)
					return
				case ch <- out.Val:
				}
			}
		}
	}()

	return ch, ctl
}

func (cl *Client) ItemLongPolling(name string) chan ItemError {
	ch := make(chan ItemError)
	outch, _ := cl.request("GET", "/items/"+name, "", &Item{}, LongPolling)
	go func() {
		it := ItemError{}
		resp := <-outch
		it.Error = resp.Error
		if resp.Error == nil {
			item, ok := resp.Val.(Item)
			if ok {
				it.Item = item
			} else {
				it.Error = fmt.Errorf("Unable to convert interface")
			}
		}
		ch <- it
		close(ch)
		return
	}()

	return ch
}

func (cl *Client) PageLongPolling(smap, page string) chan PageError {
	ch := make(chan PageError)
	outch, _ := cl.request("GET", "/sitemaps/"+smap+"/"+page, "", &SitemapPage{}, LongPolling)
	go func() {
		it := PageError{}
		resp := <-outch
		it.Error = resp.Error
		if resp.Error == nil {
			page, ok := resp.Val.(SitemapPage)
			if ok {
				it.Page = page
			} else {
				it.Error = fmt.Errorf("Unable to convert interface")
			}
		}
		ch <- it
		close(ch)
		return
	}()

	return ch
}

// Stub for long-polling item
func (cl *Client) ItemStreaming(name string) (chan ItemError, chan struct{}) {
	ch := make(chan ItemError)
	stream, ctl := cl.request("GET", "/items/"+name, "", &Item{}, Streaming)
	go func() {
		it := ItemError{}
		for resp := range stream {
			it.Error = resp.Error
			if resp.Error == nil {
				item, ok := resp.Val.(Item)
				if ok {
					it.Item = item
				} else {
					it.Error = fmt.Errorf("Unable to convert interface")
					ch <- it
					close(ch)
				}
				ch <- it
			} else {
				ch <- it
				close(ch)
			}

		}
	}()
	return ch, ctl
}

func (cl *Client) PageStreaming(smap, page string) (chan PageError, chan struct{}) {
	ch := make(chan PageError)
	stream, ctl := cl.request("GET", "/sitemaps/"+smap+"/"+page, "", &SitemapPage{}, Streaming)
	go func() {
		it := PageError{}
		for resp := range stream {
			it.Error = resp.Error
			if resp.Error == nil {
				item, ok := resp.Val.(SitemapPage)
				if ok {
					it.Page = item
				} else {
					it.Error = fmt.Errorf("Unable to convert interface")
					ch <- it
					close(ch)
				}
				ch <- it
			} else {
				ch <- it
				close(ch)
			}

		}
	}()
	return ch, ctl
}
