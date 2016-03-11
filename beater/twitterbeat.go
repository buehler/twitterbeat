package beater

import (
	"flag"
	"net/url"
	"strconv"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/buehler/twitterbeat/config"
	"github.com/buehler/twitterbeat/persistency"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"
)

var mapFile = flag.String("p", "twittermap.json", "Path to the persistency map json file")

type Twitterbeat struct {
	beatConfig *config.TwitterbeatConfig
	done       chan struct{}
	period     time.Duration
	api        *anaconda.TwitterApi
	collecting bool
	events     publisher.Client
	twitterMap *persistency.StringMap
}

// Creates beater
func New() *Twitterbeat {
	return &Twitterbeat{
		done: make(chan struct{}),
	}
}

/// *** Beater interface methods ***///

func (tb *Twitterbeat) HandleFlags(b *beat.Beat) {
	logp.Info("Handling beat flags")

	tb.twitterMap = persistency.NewStringMap()
	tb.twitterMap.Load(*mapFile)
}

func (bt *Twitterbeat) Config(b *beat.Beat) error {
	logp.Info("Loading configuration")

	// Load beater configuration
	var err error
	bt.beatConfig, err = config.NewTwitterbeatConfig()
	if err != nil {
		return err
	}

	return nil
}

func (bt *Twitterbeat) Setup(b *beat.Beat) error {
	logp.Info("Setup waitduration and api keys")

	bt.events = b.Events

	var err error
	bt.period, err = time.ParseDuration(*bt.beatConfig.Period)
	if err != nil {
		return err
	}

	anaconda.SetConsumerKey(*bt.beatConfig.Twitter.ConsumerKey)
	anaconda.SetConsumerSecret(*bt.beatConfig.Twitter.ConsumerSecret)
	bt.api = anaconda.NewTwitterApi(*bt.beatConfig.Twitter.AccessKey, *bt.beatConfig.Twitter.AccessSecret)

	return nil
}

func (bt *Twitterbeat) Run(b *beat.Beat) error {
	logp.Info("twitterbeat is running! Hit CTRL-C to stop it.")

	var err error
	ticker := time.NewTicker(bt.period)

	defer ticker.Stop()

	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
			if !bt.collecting {
				err = bt.collectTweets()
				if err != nil {
					return err
				}
			}
		}
	}
}

func (bt *Twitterbeat) Cleanup(b *beat.Beat) error {
	logp.Info("Cleanup api")

	bt.api.Close()
	return nil
}

func (bt *Twitterbeat) Stop() {
	logp.Info("Ctrl-C was hit, stopping.")

	close(bt.done)
}

func (bt *Twitterbeat) collectTweets() error {
	logp.Info("Collecting tweets")

	bt.collecting = true
	defer func() {
		bt.collecting = false
	}()

	sync, err, processed := make(chan byte), make(chan error), 0

	for _, name := range *bt.beatConfig.Twitter.Names {
		go bt.processUser(name, sync, err)
	}

	for {
		select {
		case <-sync:
			processed++
			if processed == len(*bt.beatConfig.Twitter.Names) {
				return nil
			}
		case e := <-err:
			return e
		}
	}

	return nil
}

func (bt *Twitterbeat) processUser(name string, sync chan byte, err chan error) {
	logp.Info("Collecting tweets for '%v'", name)

	var e error

	v := url.Values{}
	v.Set("screen_name", name)
	v.Set("trim_user", "true")

	if bt.twitterMap.Contains(name) {
		v.Set("since_id", bt.twitterMap.Get(name))
	}

	result, e := bt.api.GetUserTimeline(v)

	if e != nil {
		switch e.(type) {
		case anaconda.TwitterError:
			logp.Err("TwitterApi threw error: %v\nfor name: %v", e, name)
			sync <- 1
		default:
			logp.Critical("Non twitterapi error happend: %v", e)
			err <- e
		}
		return
	}

	logp.Info("Got %v tweets for '%v'", len(result), name)

	for _, tweet := range result {
		event := common.MapStr{
			"@timestamp": common.Time(time.Now()),
			"type":       "tweet",
			"screenName": name,
			"tweet":      tweet,
		}

		bt.events.PublishEvent(event)
	}

	if len(result) >= 1 {
		bt.twitterMap.Set(name, strconv.FormatInt(result[0].Id, 10))
	}

	sync <- 1
}
