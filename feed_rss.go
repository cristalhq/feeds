package feeds

import (
	"encoding/xml"
	"fmt"
	"time"
)

// rss support
// validation done according to spec here:
//    http://cyber.law.harvard.edu/rss/rss.html

// RSSFeed ...
type RSSFeed struct {
	XMLName        xml.Name `xml:"channel"`
	Title          string   `xml:"title"`       // required
	Link           string   `xml:"link"`        // required
	Description    string   `xml:"description"` // required
	Language       string   `xml:"language,omitempty"`
	Copyright      string   `xml:"copyright,omitempty"`
	ManagingEditor string   `xml:"managingEditor,omitempty"` // Author used
	WebMaster      string   `xml:"webMaster,omitempty"`
	PublishDate    string   `xml:"pubDate,omitempty"`       // created or updated
	LastBuildDate  string   `xml:"lastBuildDate,omitempty"` // updated used
	Category       string   `xml:"category,omitempty"`
	Generator      string   `xml:"generator,omitempty"`
	Docs           string   `xml:"docs,omitempty"`
	Cloud          string   `xml:"cloud,omitempty"`
	TTL            int      `xml:"ttl,omitempty"`
	Rating         string   `xml:"rating,omitempty"`
	SkipHours      string   `xml:"skipHours,omitempty"`
	SkipDays       string   `xml:"skipDays,omitempty"`
	Image          *RSSImage
	TextInput      *RSSTextInput
	Items          []*RSSItem `xml:"item"`
}

// RSSItem ...
type RSSItem struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`       // required
	Link        string   `xml:"link"`        // required
	Description string   `xml:"description"` // required
	Content     *RSSContent
	Author      string `xml:"author,omitempty"`
	Category    string `xml:"category,omitempty"`
	Comments    string `xml:"comments,omitempty"`
	Enclosure   *RSSEnclosure
	GUID        string `xml:"guid,omitempty"`    // Id used
	PubDate     string `xml:"pubDate,omitempty"` // created or updated
	Source      string `xml:"source,omitempty"`
}

// RSSFeedXML ...
// private wrapper around the RSSFeed which gives us the <rss>..</rss> xml
type RSSFeedXML struct {
	XMLName          xml.Name `xml:"rss"`
	Version          string   `xml:"version,attr"`
	ContentNamespace string   `xml:"xmlns:content,attr"`
	Channel          *RSSFeed
}

// RSSContent ...
type RSSContent struct {
	XMLName xml.Name `xml:"content:encoded"`
	Content string   `xml:",cdata"`
}

// RSSImage ...
type RSSImage struct {
	XMLName xml.Name `xml:"image"`
	URL     string   `xml:"url"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
	Width   int      `xml:"width,omitempty"`
	Height  int      `xml:"height,omitempty"`
}

// RSSTextInput ...
type RSSTextInput struct {
	XMLName     xml.Name `xml:"textInput"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Name        string   `xml:"name"`
	Link        string   `xml:"link"`
}

// RSSEnclosure ...
type RSSEnclosure struct {
	XMLName xml.Name `xml:"enclosure"`
	URL     string   `xml:"url,attr"`
	Length  string   `xml:"length,attr"`
	Type    string   `xml:"type,attr"`
}

var _ xml.Marshaler = &RSS{}

// RSS ...
type RSS struct {
	*Feed
}

// MarshalXML ...
func (r *RSS) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(r, start)
}

// RSSFeed create a new RSSFeed with a generic Feed struct's data
func (r *RSS) RSSFeed() *RSSFeed {
	var author string
	if r.Author != nil {
		author = r.Author.Email
		if len(r.Author.Name) > 0 {
			author = fmt.Sprintf("%s (%s)", r.Author.Email, r.Author.Name)
		}
	}

	var image *RSSImage
	if r.Image != nil {
		image = &RSSImage{
			URL:    r.Image.URL,
			Title:  r.Image.Title,
			Link:   r.Image.Link,
			Width:  r.Image.Width,
			Height: r.Image.Height,
		}
	}

	feed := &RSSFeed{
		Title:          r.Title,
		Link:           r.Link.Href,
		Description:    r.Description,
		ManagingEditor: author,
		PublishDate:    anyTimeFormat(time.RFC1123Z, r.Created, r.Updated),
		LastBuildDate:  anyTimeFormat(time.RFC1123Z, r.Updated),
		Copyright:      r.Copyright,
		Image:          image,
	}

	for _, i := range r.Items {
		feed.Items = append(feed.Items, newRSSItem(i))
	}
	return feed
}

// create a new RSSItem with a generic Item struct's data
func newRSSItem(i *Item) *RSSItem {
	item := &RSSItem{
		Title:       i.Title,
		Link:        i.Link.Href,
		Description: i.Description,
		GUID:        i.ID,
		PubDate:     anyTimeFormat(time.RFC1123Z, i.Created, i.Updated),
	}
	if len(i.Content) > 0 {
		item.Content = &RSSContent{Content: i.Content}
	}
	if i.Source != nil {
		item.Source = i.Source.Href
	}

	if i.Enclosure != nil && i.Enclosure.Type != "" && i.Enclosure.Length != "" {
		item.Enclosure = &RSSEnclosure{
			URL:    i.Enclosure.URL,
			Type:   i.Enclosure.Type,
			Length: i.Enclosure.Length,
		}
	}

	if i.Author != nil {
		item.Author = i.Author.Name
	}
	return item
}

// // FeedXML returns an XML-ready object for an RssFeed object
// func (r *RssFeed) FeedXML() interface{} {
// 	return &RssFeedXML{
// 		Version:          "2.0",
// 		Channel:          r,
// 		ContentNamespace: "http://purl.org/rss/1.0/modules/content/",
// 	}
// }
