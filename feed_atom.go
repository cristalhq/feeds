package feeds

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"time"
)

const xmlNamespace = "http://www.w3.org/2005/Atom"

// https://tools.ietf.org/html/rfc2026

// AtomFeed ...
type AtomFeed struct {
	XMLName     xml.Name `xml:"feed"`
	Xmlns       string   `xml:"xmlns,attr"`
	Title       string   `xml:"title"`   // required
	ID          string   `xml:"id"`      // required
	Updated     string   `xml:"updated"` // required
	Category    string   `xml:"category,omitempty"`
	Icon        string   `xml:"icon,omitempty"`
	Logo        string   `xml:"logo,omitempty"`
	Rights      string   `xml:"rights,omitempty"` // copyright used
	Subtitle    string   `xml:"subtitle,omitempty"`
	Link        *AtomLink
	Author      *AtomAuthor `xml:"author,omitempty"`
	Contributor *AtomContributor
	Entries     []*AtomEntry `xml:"entry"`
}

// AtomEntry ...
type AtomEntry struct {
	XMLName     xml.Name `xml:"entry"`
	Xmlns       string   `xml:"xmlns,attr,omitempty"`
	Title       string   `xml:"title"`   // required
	Updated     string   `xml:"updated"` // required
	ID          string   `xml:"id"`      // required
	Category    string   `xml:"category,omitempty"`
	Content     *AtomContent
	Rights      string `xml:"rights,omitempty"`
	Source      string `xml:"source,omitempty"`
	Published   string `xml:"published,omitempty"`
	Contributor *AtomContributor
	Links       []AtomLink   // required if no child 'content' elements
	Summary     *AtomSummary // required if content has src or content is base64
	Author      *AtomAuthor  // required if feed lacks an author
}

// AtomLink ...
// Multiple links with different rel can coexist
type AtomLink struct {
	XMLName  xml.Name `xml:"link"`
	Href     string   `xml:"href,attr"`
	Rel      string   `xml:"rel,attr,omitempty"`
	Type     string   `xml:"type,attr,omitempty"`
	Hreflang string   `xml:"hreflang,attr,omitempty"`
	Text     string   `xml:"text,attr,omitempty"`
	Length   string   `xml:"length,attr,omitempty"`
}

// AtomSummary ...
type AtomSummary struct {
	XMLName xml.Name `xml:"summary"`
	Content string   `xml:",chardata"`
	Type    string   `xml:"type,attr"`
}

// AtomContent ...
type AtomContent struct {
	XMLName xml.Name `xml:"content"`
	Content string   `xml:",chardata"`
	Type    string   `xml:"type,attr"`
}

// AtomAuthor ...
type AtomAuthor struct {
	XMLName xml.Name `xml:"author"`
	AtomPerson
}

// AtomContributor ...
type AtomContributor struct {
	XMLName xml.Name `xml:"contributor"`
	AtomPerson
}

// AtomPerson ...
type AtomPerson struct {
	Name  string `xml:"name,omitempty"`
	URI   string `xml:"uri,omitempty"`
	Email string `xml:"email,omitempty"`
}

var _ xml.Marshaler = &Atom{}

// Atom ...
type Atom struct {
	*Feed
}

// MarshalXML ...
func (a *Atom) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(a, start)
}

// AtomFeed create a new AtomFeed with a generic Feed struct's data
func (a *Atom) AtomFeed() *AtomFeed {
	feed := &AtomFeed{
		Xmlns: xmlNamespace,
		Title: a.Title,
		Link: &AtomLink{
			Href: a.Link.Href,
			Rel:  a.Link.Rel,
		},
		Subtitle: a.Description,
		ID:       a.Link.Href,
		Updated:  anyTimeFormat(time.RFC3339, a.Updated, a.Created),
		Rights:   a.Copyright,
	}
	if a.Author != nil {
		feed.Author = &AtomAuthor{
			AtomPerson: AtomPerson{
				Name:  a.Author.Name,
				Email: a.Author.Email,
			},
		}
	}

	feed.Entries = make([]*AtomEntry, 0, len(a.Items))
	for _, e := range a.Items {
		feed.Entries = append(feed.Entries, newAtomEntry(e))
	}
	return feed
}

func newAtomEntry(i *Item) *AtomEntry {
	id := i.ID
	// assume the description is html

	if len(id) == 0 {
		// if there's no id set, try to create one, either from data or just a uuid
		if len(i.Link.Href) > 0 && (!i.Created.IsZero() || !i.Updated.IsZero()) {
			dateStr := anyTimeFormat("2006-01-02", i.Updated, i.Created)
			host, path := i.Link.Href, "/invalid.html"
			if url, err := url.Parse(i.Link.Href); err == nil {
				host, path = url.Host, url.Path
			}
			id = fmt.Sprintf("tag:%s,%s:%s", host, dateStr, path)
		} else {
			// id = "urn:uuid:" + NewUUID().String()
		}
	}
	var name, email string
	if i.Author != nil {
		name, email = i.Author.Name, i.Author.Email
	}

	linkrel := i.Link.Rel
	if linkrel == "" {
		linkrel = "alternate"
	}
	x := &AtomEntry{
		Title:   i.Title,
		Links:   []AtomLink{{Href: i.Link.Href, Rel: linkrel, Type: i.Link.Type}},
		ID:      id,
		Updated: anyTimeFormat(time.RFC3339, i.Updated, i.Created),
		Summary: &AtomSummary{Content: i.Description, Type: "html"},
	}

	// if there's a content, assume it's html
	if len(i.Content) > 0 {
		x.Content = &AtomContent{Content: i.Content, Type: "html"}
	}

	if i.Enclosure != nil && linkrel != "enclosure" {
		x.Links = append(x.Links, AtomLink{
			Href:   i.Enclosure.URL,
			Rel:    "enclosure",
			Type:   i.Enclosure.Type,
			Length: i.Enclosure.Length},
		)
	}

	if len(name) > 0 || len(email) > 0 {
		x.Author = &AtomAuthor{AtomPerson: AtomPerson{Name: name, Email: email}}
	}
	return x
}

// // FeedXML returns an XML-Ready object for an Atom object
// func (a *Atom) FeedXML() interface{} {
// 	return a.AtomFeed()
// }

// // FeedXML returns an XML-ready object for an AtomFeed object
// func (a *AtomFeed) FeedXML() interface{} {
// 	return a
// }
