package openhab

import (
	"encoding/json"
	"encoding/xml"
)

type Link struct {
	XMLName xml.Name `xml:"link" json:"-"`
	Type    string   `json:"@type,omitempty" xml:"type,attr,omitempty"`
	URL     string   `json:"$,omitempty" xml:",chardata"`
}

type RestBase struct {
	XMLName xml.Name `xml:"openhab" json:"-"`
	Links   []Link   `json:"link,omitempty" xml:"links,omitempty"`
}

type Item struct {
	XMLName xml.Name `xml:"item" json:"-"`
	Link    string   `json:"link,omitempty" xml:"link,omitempty"`
	Name    string   `json:"name,omitempty" xml:"name,omitempty"`
	State   string   `json:"state,omitempty" xml:"state,omitempty"`
	Type    string   `json:"type,omitempty" xml:"type,omitempty"`
	Members Items    `json:"members,omitempty" xml:"members,omitempty"`
}

type Items []Item
type ItemsResp struct {
	XMLName xml.Name `xml:"items" json:"-"`
	Items   Items    `json:"item" xml:"items,omitempty"`
}

func (i *Items) UnmarshalJSON(bs []byte) error {
	single := Item{}
	multiple := make([]Item, 0)
	err := json.Unmarshal(bs, &multiple)
	if err != nil {
		err := json.Unmarshal(bs, &single)
		if err != nil {
			return err
		}
		*i = []Item{single}
		return nil
	}
	*i = multiple
	return nil
}

type Sitemaps []Sitemap
type SitemapsResp struct {
	XMLName  xml.Name `xml:"sitemaps" json:"-"`
	Sitemaps Sitemaps `json:"sitemap,omitempty" xml:"sitemap,omitempty"`
}

func (s *Sitemaps) UnmarshalJSON(bs []byte) error {
	single := Sitemap{}
	multiple := make([]Sitemap, 0)
	err := json.Unmarshal(bs, &multiple)
	if err != nil {
		err := json.Unmarshal(bs, &single)
		if err != nil {
			return err
		}
		*s = []Sitemap{single}
		return nil
	}
	*s = multiple
	return nil
}

type Sitemap struct {
	XMLName  xml.Name     `xml:"sitemap" json:"-"`
	Homepage *SitemapPage `json:"homepage,omitempty" xml:"homepage,omitempty"`
	Label    string       `json:"label,omitempty" xml:"label,omitempty"`
	Link     string       `json:"link,omitempty" xml:"link,omitempty"`
	Name     string       `json:"name,omitempty" xml:"name,omitempty"`
}

type Widget struct {
	XMLName       xml.Name     `xml:"widget" json:"-"`
	Icon          string       `json:"icon,omitempty" xml:"icon,omitempty"`
	Item          *Item        `json:"item,omitempty" xml:"item,omitempty"`
	Label         string       `json:"label,omitempty" xml:"label,omitempty"`
	Type          string       `json:"type,omitempty" xml:"type,omitempty"`
	WidgetId      string       `json:"widgetId,omitempty" xml:"widgetId,omitempty"`
	LinkedPage    *SitemapPage `json:"linkedPage,omitempty" xml:"linkedPage,omitempty"`
	SendFrequency string       `json:"sendFrequency,omitempty" xml:"sendFrequency,omitempty"`
	SwitchSupport string       `json:"switchSupport,omitempty" xml:"switchSupport,omitempty"`
	Mappings      Mappings     `json:"mapping,omitempty" xml:"mappings,omitempty"`
}

type Mappings []Mapping

func (m *Mappings) UnmarshalJSON(bs []byte) error {
	single := Mapping{}
	multiple := make([]Mapping, 0)
	err := json.Unmarshal(bs, &multiple)
	if err != nil {
		err := json.Unmarshal(bs, &single)
		if err != nil {
			return err
		}
		*m = []Mapping{single}
		return nil
	}
	*m = multiple
	return nil
}

type Mapping struct {
	XMLName xml.Name `xml:"mapping" json:"-"`
	Command string   `json:"command,omitempty" xml:"command,omitempty"`
	Label   string   `json:"label,omitempty" xml:"label,omitempty"`
}

type SitemapPage struct {
	Icon    string  `json:"icon,omitempty" xml:"icon,omitempty,omitempty"`
	Id      string  `json:"id,omitempty" xml:"id,omitempty,omitempty"`
	Leaf    string  `json:"leaf,omitempty" xml:"leaf,omitempty,omitempty"`
	Link    string  `json:"link,omitempty" xml:"link,omitempty,omitempty"`
	Title   string  `json:"title,omitempty" xml:"title,omitempty,omitempty"`
	Widgets Widgets `json:"widget,omitempty" xml:"widgets,omitempty,omitempty"`
}

type Widgets []Widget

func (w *Widgets) UnmarshalJSON(bs []byte) error {
	single := Widget{}
	multiple := make([]Widget, 0)
	err := json.Unmarshal(bs, &multiple)
	if err != nil {
		err := json.Unmarshal(bs, &single)
		if err != nil {
			return err
		}
		*w = []Widget{single}
		return nil
	}
	*w = multiple
	return nil
}
