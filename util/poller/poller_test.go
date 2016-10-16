package poller

import (
	"testing"
	"time"

	lpkg "github.com/pestophagous/hackybeat/util/logger"
)

type simplePollable struct {
	counter int
}

func (this *simplePollable) DoPoll() time.Duration {
	this.counter++
	return 0
}

func (this *simplePollable) OnShutdown() {
}

func TestPollBasic(t *testing.T) {

	logger := lpkg.NewNoopLogAdapter()
	pollable := &simplePollable{}

	poller := NewPoller(logger, pollable)

	poller.BeginBackgroundPolling()

	poller.Stop()

	if pollable.counter < 1 {
		t.Error("Pollable was not polled.")
	}
}
