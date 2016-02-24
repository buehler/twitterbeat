package beater

import (
	"fmt"

	"flag"
	"github.com/ChimeraCoder/anaconda"
	"github.com/buehler/PersistentStringMap/persistency"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/cfgfile"
	"github.com/elastic/beats/libbeat/logp"
	"time"
)

var mapFile = flag.String("p", "twittermap.json", "Path to the persistency map json file")

type TwitterBeat struct {
	alive      chan byte
	api        *anaconda.TwitterApi
	config     *TwitterConfig
	period     time.Duration
	twitterMap *persistency.PersistentStringMap
}

func New() *TwitterBeat {
	beat := &TwitterBeat{}
	return beat
}

func (tb *TwitterBeat) HandleFlags(b *beat.Beat) {
	tb.twitterMap = persistency.NewPersistentStringMap()
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

	anaconda.SetConsumerKey(*tb.config.Twitter.ConsumerKey)
	anaconda.SetConsumerSecret(*tb.config.Twitter.ConsumerSecret)
	tb.api = anaconda.NewTwitterApi(*tb.config.Twitter.AccessKey, *tb.config.Twitter.AccessSecret)

	return nil
}

func (tb *TwitterBeat) Run(b *beat.Beat) error {
	//var err error
	ticker := time.NewTicker(tb.period)

	defer ticker.Stop()

	for {
		select {
		case <-tb.alive:
			return nil
		case <-ticker.C:
			fmt.Println("Run")

		}
	}

	return nil
}

func (tb *TwitterBeat) Cleanup(b *beat.Beat) error {
	fmt.Println("Cleanup")
	return nil
}

func (tb *TwitterBeat) Stop() {
	fmt.Println("Stop")
	close(tb.alive)
}
