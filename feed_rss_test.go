package feeds

import (
	"encoding/xml"
	"testing"
)

func Test1(t *testing.T) {
	v := &RssFeedXML{
		Version: "lel",
	}

	raw, err := xml.Marshal(v)

	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(raw))
}
