package rsspoll

import (
	"time"

	"github.com/elastic/beats/libbeat/common"
	rss "github.com/jteeuwen/go-pkg-rss"

	pollcommon "github.com/pestophagous/hackybeat/beats-pollables/common"
	"github.com/pestophagous/hackybeat/util/poller"
)

func init() {
	p := poller.NewPoller(pollcommon.Logger, newPolledFeed(pollcommon.Logger, receiveRssItem))
	pollcommon.RegisterPoller(p)
}

// type polledFeed struct calls here when an Item is ready. this method converts and forwards the item to libbeat publisher
func receiveRssItem(item *rss.Item) {

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

	pollcommon.BeatsPublish(event)
}