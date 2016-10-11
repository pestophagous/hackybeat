package poller

import (
	"sync"
	"time"

	lpkg "github.com/pestophagous/hackybeat/logger"
)

type Pollable interface {
	DoPoll() time.Duration
	OnShutdown()
}

// ---------------------------------------------------------

const briefestAllowedGap = time.Second * 60

type Poller struct {
	stopperChan chan bool
	waitGroup   *sync.WaitGroup
	poll        Pollable
	logger      *lpkg.LogWithNilCheck
}

func NewPoller(log *lpkg.LogAdapter, toBePolled Pollable) *Poller {
	p := &Poller{
		stopperChan: make(chan bool),
		waitGroup:   &sync.WaitGroup{},
		poll:        toBePolled,
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