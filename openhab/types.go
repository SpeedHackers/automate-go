package openhab

import "encoding/json"

type Link struct {
	Type string `json:"@type,omitempty"`
	URL  string `json:"$,omitempty"`
}

type RestBase struct {
	Links []Link `json:"link,omitempty"`
}

type Item struct {
	Link  string `json:"link,omitempty"`
	Name  string `json:"name,omitempty"`
	State string `json:"state,omitempty"`
	Type  string `json:"type,omitempty"`
}

type Items []Item
type ItemsResp struct {
	Items Items `json:"item"`
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
	Sitemaps Sitemaps `json:"sitemap,omitempty"`
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
	Homepage *SitemapPage `json:"homepage,omitempty"`
	Label    string       `json:"label,omitempty"`
	Link     string       `json:"link,omitempty"`
	Name     string       `json:"name,omitempty"`
}

type Widget struct {
	Icon          string       `json:"icon,omitempty"`
	Item          *Item        `json:"item,omitempty"`
	Label         string       `json:"label,omitempty"`
	Type          string       `json:"type,omitempty"`
	WidgetId      string       `json:"widgetId,omitempty"`
	LinkedPage    *SitemapPage `json:"linkedPage,omitempty"`
	SendFrequency string       `json:"sendFrequency,omitempty"`
	SwitchSupport string       `json:"switchSupport,omitempty"`
	Mappings      Mappings     `json:"mapping,omitempty"`
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
	Command string `json:"command,omitempty"`
	Label   string `json:"label,omitempty"`
}

type SitemapPage struct {
	Icon    string  `json:"icon,omitempty"`
	Id      string  `json:"id,omitempty"`
	Leaf    string  `json:"leaf,omitempty"`
	Link    string  `json:"link,omitempty"`
	Title   string  `json:"title,omitempty"`
	Widgets Widgets `json:"widget,omitempty"`
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
