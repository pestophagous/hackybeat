package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/pestophagous/hackybeat/config"
	"github.com/pestophagous/hackybeat/rss-poll"
)

type Hackybeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
	poller *rsspoll.Poller
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

	bt.poller = rsspoll.NewPoller(&rsspoll.LogAdapter{
		Err:  logp.Err,
		Warn: logp.Info,
		Info: logp.Info,
		Debug: func(format string, v ...interface{}) {
			logp.Debug("hackybeat", format, v)
		},
	})
	bt.poller.BeginBackgroundPolling()

	bt.client = b.Publisher.Connect()
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
