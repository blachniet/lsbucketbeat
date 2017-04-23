// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period  time.Duration `config:"period"`
	Buckets []Bucket      `config:"buckets"`
}

type Bucket struct {
	Title       string        `config:"title"`
	Dir         string        `config:"dir"`
	FilePattern string        `config:"filePattern"`
	RetryCount  int           `config:"retryCount"`
	RetryDelay  time.Duration `config:"retryDelay"`
}

var DefaultConfig = Config{
	Period:  1 * time.Minute,
	Buckets: []Bucket{},
}
