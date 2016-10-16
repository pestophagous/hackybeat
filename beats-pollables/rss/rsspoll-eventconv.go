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

func choosyAppend(strings []string, s string) []string {
	// maybe later we'll do more sophisticated filtering using trimming or removal of 'gremlin' chars
	result := strings
	if len(s) > 0 {
		result = append(result, s)
	}

	return result
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
		categories = choosyAppend(categories, c.Domain) // <-- apparently always an empty string, but we'll check it just in case!
		categories = choosyAppend(categories, c.Text)   // <-- definitely known to contain a meaningful string value
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
