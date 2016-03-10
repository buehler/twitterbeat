package beater

import (
	"flag"
	"net/url"
	"strconv"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/buehler/go-elastic-twitterbeat/persistency"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/cfgfile"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"
)

var mapFile = flag.String("p", "twittermap.json", "Path to the persistency map json file")

type TwitterBeat struct {
	alive      chan byte
	api        *anaconda.TwitterApi
	collecting bool
	config     *TwitterConfig
	events     publisher.Client
	period     time.Duration
	twitterMap *persistency.StringMap
}

func New() *TwitterBeat {
	beat := &TwitterBeat{}
	return beat
}

func (tb *TwitterBeat) HandleFlags(b *beat.Beat) {
	tb.twitterMap = persistency.NewStringMap()
	tb.twitterMap.Load(*mapFile)
}

func (tb *TwitterBeat) Config(b *beat.Beat) error {
	config := TwitterConfigYaml{}
	err := cfgfile.Read(&config, "")
	if err != nil {
		logp.Err("Error reading configuration file: %v", err)
		return err
	}
	tb.config = &config.Input

	if config.Input.Period != nil {
		tb.period = time.Duration(*config.Input.Period) * time.Second
	} else {
		tb.period = 60 * time.Second
	}

	return nil
}

func (tb *TwitterBeat) Setup(b *beat.Beat) error {
	tb.alive = make(chan byte)
	tb.events = b.Events

	anaconda.SetConsumerKey(*tb.config.Twitter.ConsumerKey)
	anaconda.SetConsumerSecret(*tb.config.Twitter.ConsumerSecret)
	tb.api = anaconda.NewTwitterApi(*tb.config.Twitter.AccessKey, *tb.config.Twitter.AccessSecret)

	return nil
}

func (tb *TwitterBeat) Run(b *beat.Beat) error {
	var err error
	ticker := time.NewTicker(tb.period)

	defer ticker.Stop()

	for {
		select {
		case <-tb.alive:
			return nil
		case <-ticker.C:
			if !tb.collecting {
				err = tb.collectTweets()
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (tb *TwitterBeat) Stop() {
	close(tb.alive)
}

func (tb *TwitterBeat) Cleanup(b *beat.Beat) error {
	tb.api.Close()
	return nil
}

func (tb *TwitterBeat) collectTweets() error {
	tb.collecting = true
	defer func() {
		tb.collecting = false
	}()

	sync, err, processed := make(chan byte), make(chan error), 0

	for _, name := range *tb.config.Twitter.Names {
		go tb.processUser(name, sync, err)
	}

	for {
		select {
		case <-sync:
			processed++
			if processed == len(*tb.config.Twitter.Names) {
				return nil
			}
		case e := <-err:
			return e
		}
	}

	return nil
}

func (tb *TwitterBeat) processUser(name string, sync chan byte, err chan error) {
	var e error

	v := url.Values{}
	v.Set("screen_name", name)
	v.Set("trim_user", "true")

	if tb.twitterMap.Contains(name) {
		v.Set("since_id", tb.twitterMap.Get(name))
	}

	result, e := tb.api.GetUserTimeline(v)

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

	for _, tweet := range result {
		event := common.MapStr{
			"@timestamp": common.Time(time.Now()),
			"type":       "tweet",
			"screenName": name,
			"tweet":      tweet,
		}

		tb.events.PublishEvent(event)
	}

	if len(result) >= 1 {
		tb.twitterMap.Set(name, strconv.FormatInt(result[0].Id, 10))
	}

	sync <- 1
}
