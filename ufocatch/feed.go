package ufocatch

import (
	"encoding/xml"
	"net/url"
	"time"
)

// Feed for ufocatch atom
type Feed struct {
	XMLName xml.Name
	Links   []Link    `xml:"link"`
	Entries []Entry   `xml:"entry"`
	Updated time.Time `xml:"updated"`
}

// Link find first Link matching ref
func (f Feed) Link(rel string) *url.URL {
	for _, l := range f.Links {
		if l.Rel == rel {
			return l.URL
		}
	}
	return nil
}

// Link holds link tag and attributes by <link> tag
type Link struct {
	Rel  string
	Type string
	URL  *url.URL
}

// UnmarshalXML implements xml.Unmarshaler interface
func (l *Link) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var ll struct {
		Rel  string `xml:"rel,attr"`
		Type string `xml:"type,attr"`
		URL  string `xml:"href,attr"`
	}

	if err := d.DecodeElement(&ll, &start); err != nil {
		return err
	}

	u, err := url.Parse(ll.URL)
	if err != nil {
		return err
	}

	*l = Link{Rel: ll.Rel, Type: ll.Type, URL: u}
	return nil
}

// Entry holds each resource data by <entry> tag
type Entry struct {
	Title   string    `xml:"title"`
	ID      string    `xml:"id"`
	DocID   string    `xml:"docid"`
	Links   []Link    `xml:"link"`
	Updated time.Time `xml:"updated"`
}

// Link searches Links by ref and format(type attribute).
// If you omit rel or format, just match every rel and format.
func (e Entry) Link(rel string, format string) []*url.URL {
	var links []*url.URL
	for _, l := range e.Links {
		if rel != "" && rel == l.Rel {
			continue
		}
		if format != "" && format == l.Type {
			continue
		}
		links = append(links, l.URL)
	}
	return links
}
