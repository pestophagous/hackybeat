package common

import (
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	dedupe "github.com/pestophagous/hackybeat/util/deduper"
	lpkg "github.com/pestophagous/hackybeat/util/logger"
	"github.com/pestophagous/hackybeat/util/poller"
)

// Logger should be the only log-sink used by any pollable. This will ensure consistent, centralized output.
var Logger *lpkg.LogWithNilCheck

var publisherFunc func(event common.MapStr, opts ...publisher.ClientOption) bool

var pollers []*poller.Poller

var deduper *dedupe.Tool

func init() {
	log := &lpkg.LogAdapter{
		Err:  logp.Err,
		Warn: logp.Info,
		Info: logp.Info,
		Debug: func(format string, v ...interface{}) {
			logp.Debug("hackybeat", format, v)
		},
	}

	Logger = &lpkg.LogWithNilCheck{log}
	deduper = dedupe.NewDeduperTool("hackybeat", Logger)
}

// Call RegisterPollable to add a poller to the application.
func RegisterPollable(pollable poller.Pollable) {
	p := poller.NewPoller(Logger, pollable)
	pollers = append(pollers, p)
}

// InstallPublisherFunc should be called by outer application code, not by the pollables. It should be called during start-up.
func InstallPublisherFunc(f func(event common.MapStr, opts ...publisher.ClientOption) bool) {
	publisherFunc = f
}

func applyToAllPollers(f func(p *poller.Poller)) {
	for _, p := range pollers {
		f(p)
	}
}

func LaunchAllPollers() {
	applyToAllPollers((*poller.Poller).BeginBackgroundPolling)
}

func StopAllPollers() {
	applyToAllPollers((*poller.Poller).Stop)
}

// BeatsPublish is intended for use by the pollables. To be called for each event detected by a pollable.
func BeatsPublish(event common.MapStr, opts ...publisher.ClientOption) {
	if publisherFunc == nil {
		panic("Must not call BeatsPublish prior to calling InstallPublisherFunc.")
	}

	var ctime common.Time = event["@timestamp"].(common.Time)
	var ttime time.Time = time.Time(ctime)
	if deduper.IsGrantingApproval(ttime, event["type"].(string), event) {
		publisherFunc(event, opts...)
	}
}
