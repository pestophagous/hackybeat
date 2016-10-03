package beater

import (
	//	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/pestophagous/hackybeat/config"
)

type Hackybeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
}

func New() *Hackybeat {
	return &Hackybeat{
		done: make(chan struct{}),
	}
}

// Creates beater
func Newx(b *beat.Beat /*, cfg *common.Config*/) (beat.Beater, error) {
	config := config.DefaultConfig
	// if err := cfg.Unpack(&config); err != nil {
	// 	return nil, fmt.Errorf("Error reading config file: %v", err)
	// }

	bt := &Hackybeat{
		done:   make(chan struct{}),
		config: config,
	}
	return bt, nil
}

func (bt *Hackybeat) Config(b *beat.Beat) error {
	//read config file

	// err := cfgfile.Read(&ab.AbConfig, "")
	// if err != nil {
	// 	logp.Err("Error reading configuration file: %v", err)
	// 	return err
	// }

	return nil
}

func (bt *Hackybeat) Setup(b *beat.Beat) error {
	//ab.events = b.Events
	//ab.done = make(chan struct{})

	return nil
}

func (bt *Hackybeat) Run(b *beat.Beat) error {
	logp.Info("hackybeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Events //Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	for {
		select {
		case <-bt.done:
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
	//bt.client.Close()
	close(bt.done)
}

func (bt *Hackybeat) Cleanup(b *beat.Beat) error {
	return nil
}
