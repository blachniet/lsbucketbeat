package beater

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"os"

	"sync"

	"github.com/blachniet/lsbucketbeat/config"
)

type Lsbucketbeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Lsbucketbeat{
		done:   make(chan struct{}),
		config: config,
	}
	return bt, nil
}

func (bt *Lsbucketbeat) Run(b *beat.Beat) error {
	logp.Info("lsbucketbeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		wg := sync.WaitGroup{}
		for _, bkt := range bt.config.Buckets {
			wg.Add(1)
			go func(bkt config.Bucket) {
				defer wg.Done()
				bt.ls(bkt, b.Name)
			}(bkt)
		}

		wg.Wait()
	}
}

func (bt *Lsbucketbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

func (bt *Lsbucketbeat) ls(bkt config.Bucket, beatname string) {
	var contents []os.FileInfo
	err := retry(func() error {
		c, e := ioutil.ReadDir(bkt.Dir)
		contents = c
		return e
	}, bkt.RetryCount, bkt.RetryDelay)

	if err != nil {
		bt.client.PublishEvent(common.MapStr{
			"@timestamp": common.Time(time.Now().UTC()),
			"type":       beatname,
			"bucket": common.MapStr{
				"title": bkt.Title,
			},
			"error": err.Error(),
		})
	} else {
		fileCount := 0
		for _, f := range contents {
			if f.IsDir() {
				continue
			}

			fileCount++
			isMatch, _ := filepath.Match(bkt.FilePattern, f.Name())
			if isMatch {
				bt.client.PublishEvent(common.MapStr{
					"@timestamp": common.Time(time.Now().UTC()),
					"type":       beatname,
					"bucket": common.MapStr{
						"title": bkt.Title,
					},
					"file": common.MapStr{
						"path":    filepath.Join(bkt.Dir, f.Name()),
						"dir":     bkt.Dir,
						"name":    f.Name(),
						"size":    f.Size(),
						"modTime": common.Time(f.ModTime()),
					},
				})
			}
		}

		bt.client.PublishEvent(common.MapStr{
			"@timestamp": common.Time(time.Now().UTC()),
			"type":       beatname,
			"bucket": common.MapStr{
				"title": bkt.Title,
			},
			"fileCount": fileCount,
		})
	}
}

func retry(op func() error, count int, delay time.Duration) error {
	var err error
	for i := 0; i <= count; i++ {
		err = op()
		if err == nil {
			return nil
		}
	}
	return err
}
