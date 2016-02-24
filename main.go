package main

import (
	"os"

	"github.com/buehler/TwitterBeat/beater"
	"github.com/elastic/beats/libbeat/beat"
)

var name = "twitterbeat"

func main() {
	if err := beat.Run(name, "", beater.New()); err != nil {
		os.Exit(1)
	}
}
