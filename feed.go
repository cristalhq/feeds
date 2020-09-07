package feeds

import (
	"sort"
	"time"
)

// Feed ...
type Feed struct {
	Title       string    `json:"title,omitempty"`
	Link        *Link     `json:"link,omitempty"`
	Description string    `json:"description,omitempty"`
	Author      *Author   `json:"author,omitempty"`
	Updated     time.Time `json:"updated,omitempty"`
	Created     time.Time `json:"created,omitempty"`
	ID          string    `json:"id,omitempty"`
	Subtitle    string    `json:"subtitle,omitempty"`
	Items       []*Item   `json:"items,omitempty"`
	Copyright   string    `json:"copyright,omitempty"`
	Image       *Image    `json:"image,omitempty"`
}

// Item ...
type Item struct {
	ID          string     `json:"id,omitempty"` // used as GUID in RSS, id in Atom
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"` // used as description in RSS, summary in Atom
	Link        *Link      `json:"link,omitempty"`
	Source      *Link      `json:"source,omitempty"`
	Author      *Author    `json:"author,omitempty"`
	Created     time.Time  `json:"created,omitempty"`
	Updated     time.Time  `json:"updated,omitempty"`
	Enclosure   *Enclosure `json:"enclosure,omitempty"`
	Content     string     `json:"content,omitempty"`
}

// Link ...
type Link struct {
	Href   string `json:"href,omitempty"`
	Rel    string `json:"rel,omitempty"`
	Type   string `json:"type,omitempty"`
	Length string `json:"length,omitempty"`
}

// Author ...
type Author struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

// Image ...
type Image struct {
	URL    string `json:"url,omitempty"`
	Title  string `json:"title,omitempty"`
	Link   string `json:"link,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

// Enclosure ...
type Enclosure struct {
	URL    string `json:"url,omitempty"`
	Length string `json:"length,omitempty"`
	Type   string `json:"type,omitempty"`
}

// Add a new item to a feed.
func (f *Feed) Add(item *Item) {
	f.Items = append(f.Items, item)
}

// Sort items with a given compare function.
func (f *Feed) Sort(less func(a, b *Item) bool) {
	sort.SliceStable(f.Items, func(i, j int) bool {
		return less(f.Items[i], f.Items[j])
	})
}

// ToAtom ... creates an Atom representation of this feed
func (f *Feed) ToAtom() *Atom {
	return &Atom{f}
}

// ToJSON creates a JSON Feed representation of this feed
func (f *Feed) ToJSON() *JSON {
	return &JSON{f}
}

// ToRSS ... creates an RSS representation of this feed
func (f *Feed) ToRSS() *RSS {
	return &RSS{f}
}

// // XMLFeed interface used by ToXML to get a object suitable for exporting XML.
// type XMLFeed interface {
// 	FeedXML() interface{}
// }

// // ToXML turn a feed object (either a Feed, AtomFeed, or RssFeed) into xml
// // returns an error if xml marshaling fails
// func ToXML(feed XMLFeed) (string, error) {
// 	x := feed.FeedXML()
// 	data, err := xml.MarshalIndent(x, "", "  ")
// 	if err != nil {
// 		return "", err
// 	}
// 	// strip empty line from default xml header
// 	s := xml.Header[:len(xml.Header)-1] + string(data)
// 	return s, nil
// }

// WriteXML writes a feed object (either a Feed, AtomFeed, or RssFeed) as XML into
// the writer. Returns an error if XML marshaling fails.
// func WriteXML(feed XMLFeed, w io.Writer) error {
// 	var xmlHeader = []byte(`<?xml version="1.0" encoding="UTF-8"?>`)

// 	x := feed.FeedXML()
// 	if _, err := w.Write(xmlHeader); err != nil {
// 		return err
// 	}
// 	e := xml.NewEncoder(w)
// 	e.Indent("", "  ")
// 	return e.Encode(x)
// }

// // WriteAtom ... WriteAtom writes an Atom representation of this feed to the writer.
// func (f *Feed) WriteAtom(w io.Writer) error {
// 	return WriteXML(&Atom{f}, w)
// }

// // WriteRss writes an RSS representation of this feed to the writer.
// func (f *Feed) WriteRss(w io.Writer) error {
// 	return WriteXML(&Rss{f}, w)
// }

// // WriteJSON writes an JSON representation of this feed to the writer.
// func (f *Feed) WriteJSON(w io.Writer) error {
// 	j := &JSON{f}
// 	feed := j.JSONFeed()

// 	e := json.NewEncoder(w)
// 	e.SetIndent("", "  ")
// 	return e.Encode(feed)
// }
