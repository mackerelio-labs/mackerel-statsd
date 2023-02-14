package statsd_test

import (
	"context"
	"log"

	"github.com/hatena/mackerelstatsd/parser"
	"github.com/hatena/mackerelstatsd/statsd"
)

func Example_server() {
	var s statsd.Server
	ctx := context.Background()
	log.Fatal(s.ListenAndServe(ctx, ":8125", func(addr string, m *parser.Metric) error {
		log.Println(addr, *m)
		return nil
	}))
}
