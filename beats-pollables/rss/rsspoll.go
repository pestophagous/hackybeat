package rsspoll

import (
	"time"

	rss "github.com/jteeuwen/go-pkg-rss"
	lpkg "github.com/pestophagous/hackybeat/util/logger"
)

const rssUri = "http://stackoverflow.com/feeds/tag?tagnames=go%20or%20goroutine%20or%20json%20or%20python%20or%20c%2b%2b%20or%20git%20or%20linux%20or%20gdb%20or%20xcode&sort=newest"

type polledFeed struct {
	feed     *rss.Feed
	logger   *lpkg.LogWithNilCheck
	callback RssItemCallback
}

func NewPolledFeed(log *lpkg.LogAdapter, conv RssItemCallback) *polledFeed {
	p := new(polledFeed)
	p.feed = rss.New(5, true, p.chanHandler, p.itemHandler)
	p.logger = &lpkg.LogWithNilCheck{log}
	p.callback = conv
	return p
}

func (this *polledFeed) chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
	this.logger.Info("%d new channel(s) in %s\n", len(newchannels), feed.Url)
}

func (this *polledFeed) itemHandler(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
	this.logger.Info("%d new item(s) in %s\n", len(newitems), feed.Url)
	for _, item := range newitems {
		this.callback.ReceiveRssItem(item)
	}
}

func (this *polledFeed) OnShutdown() {
	this.logger.ReleaseLog()
}

func (this *polledFeed) DoPoll() time.Duration {
	// when Fetch finds new information, execution will enter chanHandler and/or itemHandler
	if err := this.feed.Fetch(rssUri, nil); err != nil {
		this.logger.Err("[e] %s: %s\n", rssUri, err)
	}

	// let the feed tell us (via SecondsTillUpdate) when it thinks we should call Fetch again
	return time.Duration(this.feed.SecondsTillUpdate() * 1e9)
}
