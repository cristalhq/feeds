package feeds

import (
	"encoding/json"
	"time"
)

const jsonFeedVersion = "https://jsonfeed.org/version/1"

// JSONFeed format is a pragmatic syndication format,
// like RSS and Atom, but with one big difference: it’s JSON instead of XML.
//
// See: https://jsonfeed.org/version/1
//
type JSONFeed struct {
	Version     string      `json:"version"`
	Title       string      `json:"title"` // required
	HomePageURL string      `json:"home_page_url,omitempty"`
	FeedURL     string      `json:"feed_url,omitempty"`
	Description string      `json:"description,omitempty"`
	UserComment string      `json:"user_comment,omitempty"`
	NextURL     string      `json:"next_url,omitempty"`
	Icon        string      `json:"icon,omitempty"`
	Favicon     string      `json:"favicon,omitempty"`
	Author      *JSONAuthor `json:"author,omitempty"`
	Expired     *bool       `json:"expired,omitempty"`
	Hubs        []*JSONHub  `json:"hubs,omitempty"`
	Items       []*JSONItem `json:"items,omitempty"`
}

// JSONItem represents a single entry/post for the feed.
//
type JSONItem struct {
	ID            string           `json:"id"` // required
	URL           string           `json:"url,omitempty"`
	ExternalURL   string           `json:"external_url,omitempty"`
	Title         string           `json:"title,omitempty"`
	ContentHTML   string           `json:"content_html,omitempty"`
	ContentText   string           `json:"content_text,omitempty"`
	Summary       string           `json:"summary,omitempty"`
	Image         string           `json:"image,omitempty"`
	BannerImage   string           `json:"banner_,omitempty"`
	DatePublished *time.Time       `json:"date_published,omitempty"`
	DateModified  *time.Time       `json:"date_modified,omitempty"`
	Author        *JSONAuthor      `json:"author,omitempty"`
	Tags          []string         `json:"tags,omitempty"`
	Attachments   []JSONAttachment `json:"attachments,omitempty"`
}

// JSONAuthor describes the feed or item author.
//
type JSONAuthor struct {
	Name   string `json:"name,omitempty"`
	URL    string `json:"url,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

// JSONHub describes endpoints that can be used to subscribe to
// real-time notifications from the publisher of this feed.
//
type JSONHub struct {
	Type string `json:"type"` // required
	URL  string `json:"url"`  // required
}

// JSONAttachment represents a related resource.
// Podcasts, for instance, would include an attachment that’s an audio or video file.
//
type JSONAttachment struct {
	URL      string `json:"url"`       // required
	MIMEType string `json:"mime_type"` // required
	Title    string `json:"title,omitempty"`
	Size     int    `json:"size,omitempty"`
	Duration int    `json:"duration_in_seconds,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (f *JSONFeed) MarshalJSON() ([]byte, error) {
	return json.Marshal(f)
}

// JSON is used to convert a generic Feed to a JSONFeed.
type JSON struct {
	*Feed
}

func (f *JSON) String() string {
	b, _ := f.JSONFeed().MarshalJSON()
	return string(b)
}

// JSONFeed creates a new JSONFeed with a generic Feed struct's data.
func (f *JSON) JSONFeed() *JSONFeed {
	feed := &JSONFeed{
		Version:     jsonFeedVersion,
		Title:       f.Title,
		Description: f.Description,
	}

	if f.Link != nil {
		feed.HomePageURL = f.Link.Href
	}
	if f.Author != nil {
		feed.Author = &JSONAuthor{Name: f.Author.Name}
	}
	for _, e := range f.Items {
		feed.Items = append(feed.Items, newJSONItem(e))
	}
	return feed
}

func newJSONItem(i *Item) *JSONItem {
	item := &JSONItem{
		ID:          i.ID,
		Title:       i.Title,
		Summary:     i.Description,
		ContentHTML: i.Content,
	}

	if i.Link != nil {
		item.URL = i.Link.Href
	}
	if i.Source != nil {
		item.ExternalURL = i.Source.Href
	}
	if i.Author != nil {
		item.Author = &JSONAuthor{Name: i.Author.Name}
	}
	if !i.Created.IsZero() {
		item.DatePublished = &i.Created
	}
	if !i.Updated.IsZero() {
		item.DateModified = &i.Updated
	}
	// if i.Enclosure != nil && strings.HasPrefix(i.Enclosure.Type, "image/") {
	// 	item.Image = i.Enclosure.URL
	// }
	return item
}
