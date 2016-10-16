package common

import (
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	lpkg "github.com/pestophagous/hackybeat/util/logger"
	"github.com/pestophagous/hackybeat/util/poller"
)

var Logger *lpkg.LogAdapter

var PublisherFunc func(event common.MapStr, opts ...publisher.ClientOption) bool

var pollers []*poller.Poller

func init() {
	Logger = &lpkg.LogAdapter{
		Err:  logp.Err,
		Warn: logp.Info,
		Info: logp.Info,
		Debug: func(format string, v ...interface{}) {
			logp.Debug("hackybeat", format, v)
		},
	}
}

func RegisterPoller(p *poller.Poller) {
	pollers = append(pollers, p)
}

func ApplyToAllPollers(f func(p *poller.Poller)) {
	for _, p := range pollers {
		f(p)
	}
}

func BeatsPublish(event common.MapStr, opts ...publisher.ClientOption) {
	if PublisherFunc == nil {
		panic("Must not call BeatsPublish prior to installing PublisherFunc.")
	}

	PublisherFunc(event, opts...)
}
