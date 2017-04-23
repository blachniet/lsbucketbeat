package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/blachniet/lsbucketbeat/beater"
)

func main() {
	err := beat.Run("lsbucketbeat", "0.1.0", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
