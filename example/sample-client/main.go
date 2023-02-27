package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/mackerelio-labs/mackerel-statsd/example/driver/cactusstatsd"
	"github.com/mackerelio-labs/mackerel-statsd/statsd"
)

func main() {
	log.SetFlags(0)
	d, err := cactusstatsd.NewDriver("localhost:8125", "custom.statsd.sample")
	if err != nil {
		log.Fatal(err)
	}
	c := statsd.NewClient(d)
	defer c.Close()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		for i := 0; i < 10; i++ {
			n := rand.Int63n(6) + 1
			if err := c.Inc("dice", n, &statsd.Options{SamplingRate: 1.0}); err != nil {
				log.Fatal(err)
			}
			time.Sleep(1 * time.Second)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for i := 0; i < 10; i++ {
			start := time.Now()
			resp, err := http.Get("https://hatena.co.jp/")
			if err != nil {
				time.Sleep(2 * time.Second)
				continue
			}
			io.Copy(io.Discard, resp.Body)
			d := time.Since(start)
			if err := c.Timer("http.hatena", d, &statsd.Options{SamplingRate: 1.0}); err != nil {
				log.Fatal(err)
			}
			resp.Body.Close()
			time.Sleep(2 * time.Second)
		}
		wg.Done()
	}()

	wg.Wait()
}
