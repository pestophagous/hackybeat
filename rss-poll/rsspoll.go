package rsspoll

import (
	"sync"
	"time"

	rss "github.com/jteeuwen/go-pkg-rss"
	lpkg "github.com/pestophagous/hackybeat/logger"
)

const rssUri = "http://stackoverflow.com/feeds/tag?tagnames=go%20or%20goroutine%20or%20json%20or%20python%20or%20c%2b%2b%20or%20git%20or%20linux%20or%20gdb%20or%20xcode&sort=newest"

type polledFeed struct {
	f      *rss.Feed
	logger *lpkg.LogWithNilCheck
}

func newPolledFeed(log *lpkg.LogAdapter) *polledFeed {
	p := new(polledFeed)
	p.f = rss.New(5, true, p.chanHandler, p.itemHandler)
	p.logger = &lpkg.LogWithNilCheck{log}
	return p
}

func (this *polledFeed) chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
	this.logger.Info("%d new channel(s) in %s\n", len(newchannels), feed.Url)
}

func (this *polledFeed) itemHandler(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
	this.logger.Info("%d new item(s) in %s\n", len(newitems), feed.Url)
}

// ---------------------------------------------------------

type Poller struct {
	stopperChan chan bool
	waitGroup   *sync.WaitGroup
	pf          *polledFeed
	logger      *lpkg.LogWithNilCheck
}

func NewPoller(log *lpkg.LogAdapter) *Poller {
	p := &Poller{
		stopperChan: make(chan bool),
		waitGroup:   &sync.WaitGroup{},
		pf:          newPolledFeed(log),
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
				// when Fetch finds new information, execution will enter chanHandler and/or itemHandler
				if err := this.pf.f.Fetch(rssUri, nil); err != nil {
					this.logger.Err("[e] %s: %s\n", rssUri, err)
				}

				// let the feed tell us (via SecondsTillUpdate) when it thinks we should call Fetch again
				if next := this.pf.f.SecondsTillUpdate() * 1e9; next > 60*1e9 {
					this.logger.Debug("%v", next)
					timer = time.NewTicker(time.Duration(next))
				} else {
					this.logger.Debug("sixty")
					timer = time.NewTicker(time.Second * 60)
				}
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
	this.pf.logger.ReleaseLog()
	this.logger.ReleaseLog()
}
