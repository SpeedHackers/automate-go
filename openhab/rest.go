package openhab

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
)

type RestError struct {
	Text string
	Code int
}

func (r RestError) Error() string {
	return fmt.Sprintf("%d: %s", r.Code, r.Text)
}

func NewRestError(err error) RestError {
	var code int
	if err == nil || len(err.Error()) < 3 {
		code = 500
	} else {
		var err2 error
		code, err2 = strconv.Atoi(err.Error()[:3])
		if err2 != nil {
			code = 500
		}
	}
	text := err.Error()
	return RestError{text, code}
}

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

type response struct {
	Val   interface{}
	Error error
}

// Needs to be updated for long-polling/streaming reqs (add a chan?)
func (cl *Client) request(method, url, body string, out interface{}, reqType ReqType) (chan response, chan struct{}) {
	ctl := make(chan struct{})
	var req *http.Request
	var err error
	ch := make(chan response, 1)
	if body != "" {
		bodyBuffer := bytes.NewBuffer([]byte(body))
		req, err = http.NewRequest(method, cl.URL+url, bodyBuffer)
		if err != nil {
			ch <- response{nil, NewRestError(err)}
			close(ch)
			return ch, nil
		}
	} else {
		req, err = http.NewRequest(method, cl.URL+url, nil)
		if err != nil {
			ch <- response{nil, NewRestError(err)}
			close(ch)
			return ch, nil
		}
	}
	req.Header.Add("Content-Type", "text/plain")
	req.Header.Add("Accept", "application/json")
	//stream := false
	switch reqType {
	case LongPolling:
		req.Header.Add("X-Atmosphere-Transport", "long-polling")
	case Streaming:
		req.Header.Add("X-Atmosphere-Transport", "streaming")
		//stream = true
	default:
	}
	if cl.Username != "" && cl.Password != "" {
		req.SetBasicAuth(cl.Username, cl.Password)
	}
	resp, err := cl.httpClient.Do(req)
	if err != nil {
		ch <- response{nil, NewRestError(err)}
		close(ch)
		return ch, nil
	}
	if resp.Status[0] == '2' {
		decoder := json.NewDecoder(resp.Body)
		if out != nil {
			dch := make(chan interface{})
			go func() {
				for {
					err := decoder.Decode(out)
					if err != nil {
						dch <- response{nil, NewRestError(err)}
						close(dch)
						return
					}
					dch <- response{out, nil}
				}
			}()
			go func() {
				for {
					select {
					case <-ctl:
						close(ch)
						return
					case <-dch:
						if err != nil {
							select {
							case <-ctl:
							case ch <- response{nil, NewRestError(err)}:
							}
							close(ch)
							return
						}
						v := reflect.ValueOf(out)
						select {
						case <-ctl:
							close(ch)
							return
						case ch <- response{reflect.Indirect(v).Interface(), nil}:
						}
					}
				}
			}()
			return ch, ctl
		}
		ch <- response{nil, nil}
		close(ch)
		return ch, nil
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ch <- response{nil, NewRestError(err)}
		close(ch)
		return ch, nil
	}
	code, _ := strconv.Atoi(resp.Status[:3])
	text := string(bodyBytes)
	err = RestError{text, code}
	ch <- response{nil, err}
	close(ch)
	return ch, nil
}
