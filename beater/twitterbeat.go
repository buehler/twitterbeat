package beater

import (
	"fmt"

	"github.com/elastic/beats/libbeat/beat"
	"time"
)

type TwitterBeat struct {
	alive  chan byte
	config TwitterConfig
}

func New() *TwitterBeat {
	return &TwitterBeat{}
}

func (tb *TwitterBeat) Config(beat *beat.Beat) error {
	fmt.Println("Config")

	return nil
}

func (tb *TwitterBeat) Setup(beat *beat.Beat) error {
	fmt.Println("Setup")
	tb.alive = make(chan byte)
	return nil
}

func (tb *TwitterBeat) Run(beat *beat.Beat) error {
	//var err error
	ticker := time.NewTicker(time.Duration(1) * time.Second)

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

func (tb *TwitterBeat) Cleanup(beat *beat.Beat) error {
	fmt.Println("Cleanup")
	return nil
}

func (tb *TwitterBeat) Stop() {
	fmt.Println("Stop")
	close(tb.alive)
}
