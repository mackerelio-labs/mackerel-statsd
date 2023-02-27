package statsd_test

import (
	"context"
	"log"

	"github.com/mackerelio-labs/mackerel-statsd/parser"
	"github.com/mackerelio-labs/mackerel-statsd/statsd"
)

func Example_server() {
	var s statsd.Server
	ctx := context.Background()
	log.Fatal(s.ListenAndServe(ctx, ":8125", func(addr string, m *parser.Metric) error {
		log.Println(addr, *m)
		return nil
	}))
}
