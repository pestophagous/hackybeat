package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	rsspoll "github.com/pestophagous/hackybeat/beats-pollables/rss"
	"github.com/pestophagous/hackybeat/config"
	lpkg "github.com/pestophagous/hackybeat/util/logger"
	"github.com/pestophagous/hackybeat/util/poller"
)

type Hackybeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
	poller *poller.Poller
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	logp.Debug("hackybeat", "New Hackybeat")
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Hackybeat{
		done:   make(chan struct{}),
		config: config,
	}
	return bt, nil
}

func (bt *Hackybeat) Run(b *beat.Beat) error {
	logp.Info("hackybeat is running! Hit CTRL-C to stop it.")

	loga := &lpkg.LogAdapter{
		Err:  logp.Err,
		Warn: logp.Info,
		Info: logp.Info,
		Debug: func(format string, v ...interface{}) {
			logp.Debug("hackybeat", format, v)
		},
	}

	bt.client = b.Publisher.Connect()

	bt.poller = poller.NewPoller(loga, rsspoll.NewPolledFeed(loga, &rsspoll.RssItemToBeatEvent{DoPublish: bt.client.PublishEvent}))

	bt.poller.BeginBackgroundPolling()

	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	for {
		select {
		case <-bt.done:
			logp.Debug("hackybeat", "case <-bt.done")
			return nil
		case <-ticker.C:
		}

		event := common.MapStr{
			"@timestamp": common.Time(time.Now()),
			"type":       b.Name,
			"counter":    counter,
		}
		bt.client.PublishEvent(event)
		logp.Info("Event sent")
		counter++
	}
}

func (bt *Hackybeat) Stop() {
	bt.poller.Stop()
	logp.Debug("hackybeat", "Stop Hackybeat")
	bt.client.Close()
	close(bt.done)
}
