package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
)

type Encoding func(io.Writer, interface{}) error

func encodeJSON(w io.Writer, i interface{}) error {
	if i == nil {
		return nil
	}
	return json.NewEncoder(w).Encode(i)
}

func encodeXML(w io.Writer, i interface{}) error {
	if i == nil {
		return nil
	}
	buf := &bytes.Buffer{}
	err := xml.NewEncoder(buf).Encode(i)
	if err != nil {
		return err
	}
	w.Write([]byte(xml.Header))
	_, err = io.Copy(w, buf)
	return err
}

var (
	XMLEncoding  Encoding = encodeXML
	JSONEncoding Encoding = encodeJSON
)
