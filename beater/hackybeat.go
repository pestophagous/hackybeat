package beater

import (
	"fmt"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	_ "github.com/pestophagous/hackybeat/beats-pollables/register"

	pollcommon "github.com/pestophagous/hackybeat/beats-pollables/common"
	"github.com/pestophagous/hackybeat/config"
	"github.com/pestophagous/hackybeat/util/poller"
)

type Hackybeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
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

	bt.client = b.Publisher.Connect()

	pollcommon.PublisherFunc = bt.client.PublishEvent

	pollcommon.ApplyToAllPollers((*poller.Poller).BeginBackgroundPolling)

	select {
	case <-bt.done:
		logp.Debug("hackybeat", "case <-bt.done")
		return nil
	}
}

func (bt *Hackybeat) Stop() {
	pollcommon.ApplyToAllPollers((*poller.Poller).Stop)
	logp.Debug("hackybeat", "Stop Hackybeat")
	bt.client.Close()
	close(bt.done)
}
