package openhab

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

func (cl *Client) longPollPath(path string, t interface{}) chan interface{} {
	ch := make(chan interface{})
	go func() {
		outch, _ := cl.request("GET", path, "", t, LongPolling)
		out := <-outch
		if out.Error != nil {
			close(ch)
			return
		}
		ch <- out.Val
	}()

	return ch
}

// Stub for long-polling item
func (cl *Client) ItemStreaming(name string) (chan Item, chan struct{}) {
	ch := make(chan Item)
	stream, ctl := cl.streamPath("/items/"+name, &Item{})
	go func() {
		for v := range stream {
			ch <- v.(Item)
		}
		close(ch)
	}()
	return ch, ctl
}

// Stub for long-polling item
func (cl *Client) PageStreaming(smap, name string) (chan SitemapPage, chan struct{}) {
	ch := make(chan SitemapPage)
	stream, ctl := cl.streamPath("/sitemaps/"+smap+"/"+name, &SitemapPage{})
	go func() {
		for v := range stream {
			ch <- v.(SitemapPage)
		}
		close(ch)
	}()
	return ch, ctl
}

// Stub for long-polling item
func (cl *Client) SitemapStreaming(smap string) (chan Sitemap, chan struct{}) {
	ch := make(chan Sitemap)
	stream, ctl := cl.streamPath("/sitemaps/"+smap, &Sitemap{})
	go func() {
		for v := range stream {
			ch <- v.(Sitemap)
		}
		close(ch)
	}()
	return ch, ctl
}

// Create a channel to receive a new item on asynchronously
func (cl *Client) ItemLongPolling(name string) chan Item {
	ch := make(chan Item)
	go func() {
		poll := cl.longPollPath("/items/"+name, &Item{})
		v := <-poll
		ch <- v.(Item)
	}()

	return ch
}

// Create a channel to receive a new item on asynchronously
func (cl *Client) SitemapLongPolling(name string) chan Sitemap {
	ch := make(chan Sitemap)
	go func() {
		poll := cl.longPollPath("/sitemaps/"+name, &Sitemap{})
		v := <-poll
		ch <- v.(Sitemap)
	}()

	return ch
}

// Create a channel to receive a new item on asynchronously
func (cl *Client) PageLongPolling(smap, name string) chan SitemapPage {
	ch := make(chan SitemapPage)
	go func() {
		poll := cl.longPollPath("/sitemaps/"+smap+"/"+name, &SitemapPage{})
		v := <-poll
		ch <- v.(SitemapPage)
	}()

	return ch
}
