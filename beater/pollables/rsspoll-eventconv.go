package pollables

import (
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/publisher"
	rss "github.com/jteeuwen/go-pkg-rss"
)

// Benefit of this interface: pollables.polledFeed does not know about types in libbeat
type RssItemCallback interface {
	ReceiveRssItem(item *rss.Item)
}

type RssItemToBeatEvent struct {
	DoPublish func(event common.MapStr, opts ...publisher.ClientOption) bool
}

// pollables.polledFeed calls here when an Item is ready. this method converts and forwards the item to libbeat publisher
func (this *RssItemToBeatEvent) ReceiveRssItem(item *rss.Item) {

	var pubDate time.Time
	pubDate, err := item.ParsedPubDate()
	if err != nil {
		pubDate = time.Now()
	}

	var categories []string
	for _, c := range item.Categories {
		categories = append(categories, c.Domain) // <-- apparently always an empty string, but i'll use it just in case!
		categories = append(categories, c.Text)   // <-- definitely known to contain a meaningful string value
	}

	// At a minimum, the event object must contain a @timestamp field and a type field. Beyond that, events can contain
	// any additional fields, and they can be created as often as necessary.
	event := common.MapStr{
		"@timestamp": common.Time(pubDate),
		"type":       "hackybeat-stackoverflow-testing",
		"title":      item.Title,
		"author":     item.Author,
		"categories": categories,
	}
	this.DoPublish(event)
}
