package rsspoll

import (
	"sync"
	"time"

	rss "github.com/jteeuwen/go-pkg-rss"
	lpkg "github.com/pestophagous/hackybeat/logger"
)

const rssUri = "http://stackoverflow.com/feeds/tag?tagnames=go%20or%20goroutine%20or%20json%20or%20python%20or%20c%2b%2b%20or%20git%20or%20linux%20or%20gdb%20or%20xcode&sort=newest"

type Pollable interface {
	DoPoll() time.Duration
	OnShutdown()
}

type polledFeed struct {
	feed   *rss.Feed
	logger *lpkg.LogWithNilCheck
}

func newPolledFeed(log *lpkg.LogAdapter) *polledFeed {
	p := new(polledFeed)
	p.feed = rss.New(5, true, p.chanHandler, p.itemHandler)
	p.logger = &lpkg.LogWithNilCheck{log}
	return p
}

func (this *polledFeed) chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
	this.logger.Info("%d new channel(s) in %s\n", len(newchannels), feed.Url)
}

func (this *polledFeed) itemHandler(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
	this.logger.Info("%d new item(s) in %s\n", len(newitems), feed.Url)
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

// ---------------------------------------------------------

const briefestAllowedGap = time.Second * 60

type Poller struct {
	stopperChan chan bool
	waitGroup   *sync.WaitGroup
	poll        Pollable
	logger      *lpkg.LogWithNilCheck
}

func NewPoller(log *lpkg.LogAdapter) *Poller {
	p := &Poller{
		stopperChan: make(chan bool),
		waitGroup:   &sync.WaitGroup{},
		poll:        newPolledFeed(log),
		logger:      &lpkg.LogWithNilCheck{log},
	}

	return p
}

func (this *Poller) BeginBackgroundPolling() {
	this.launchPolling()
}

func (this *Poller) launchPolling() {
	this.waitGroup.Add(1)
	go func() {
		defer this.waitGroup.Done()

		timer := time.NewTicker(time.Millisecond) // <-- set to 'virtually nothing' the first time through

		for {
			select {
			case <-this.stopperChan:
				return
			case <-timer.C:
				var nextAt time.Duration = this.poll.DoPoll()

				if nextAt < briefestAllowedGap {
					// log what the 'offending' nextAt was before we overwrite it
					this.logger.Info("Refusing pollable's short interval of: %s. Instead will use: %s.", nextAt, briefestAllowedGap)
					nextAt = briefestAllowedGap
				}

				timer = time.NewTicker(nextAt)
			}
		}
	}()
}

// Stop the poller by closing the poller's channel.  Block until the poller
// is really stopped.
func (this *Poller) Stop() {
	close(this.stopperChan)
	this.waitGroup.Wait()
	// func objects on LogAdapter may hold references to foreign code. Release the refs:
	this.poll.OnShutdown()
	this.logger.ReleaseLog()
}
