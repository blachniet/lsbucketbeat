package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/blachniet/lsbucketbeat/beater"
)

func main() {
	err := beat.Run("lsbucketbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
