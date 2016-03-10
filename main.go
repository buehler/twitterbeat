package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/buehler/twitterbeat/beater"
)

func main() {
	err := beat.Run("twitterbeat", "", beater.New())
	if err != nil {
		os.Exit(1)
	}
}
