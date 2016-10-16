package rsspoll

import (
	"fmt"
	"testing"

	rss "github.com/jteeuwen/go-pkg-rss"

	lpkg "github.com/pestophagous/hackybeat/util/logger"
)

func onItem(item *rss.Item) {
}

func TestRssPollBasic(t *testing.T) {

	var errorMsgs []string
	logger := lpkg.NewNoopLogAdapter()
	logger.Err = func(format string, v ...interface{}) {
		errorMsgs = append(errorMsgs, fmt.Sprintf(format, v...))
	}

	rssPoller := newPolledFeed(logger, onItem)
	rssPoller.DoPoll()

	if len(errorMsgs) > 0 {
		t.Errorf("Errors while invoking polledFeed.DoPoll: %v", errorMsgs)
	}
}
