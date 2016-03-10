package beater

import (
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/buehler/twitterbeat/config"
)

type Twitterbeat struct {
	beatConfig *config.TwitterbeatConfig
	done       chan struct{}
	period     time.Duration
}

// Creates beater
func New() *Twitterbeat {
	return &Twitterbeat{
		done: make(chan struct{}),
	}
}

/// *** Beater interface methods ***///

func (bt *Twitterbeat) Config(b *beat.Beat) error {

	// Load beater configuration
	var err error
	bt.beatConfig, err = config.NewTwitterbeatConfig()
	if err != nil {
		return err
	}

	return nil
}

func (bt *Twitterbeat) Setup(b *beat.Beat) error {

	var err error
	bt.period, err = time.ParseDuration(*bt.beatConfig.Period)
	if err != nil {
		return err
	}

	return nil
}

func (bt *Twitterbeat) Run(b *beat.Beat) error {
	logp.Info("twitterbeat is running! Hit CTRL-C to stop it.")

	ticker := time.NewTicker(bt.period)
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
		b.Events.PublishEvent(event)
		logp.Info("Event sent")
		counter++
	}
}

func (bt *Twitterbeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (bt *Twitterbeat) Stop() {
	close(bt.done)
}
