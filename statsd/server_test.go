package statsd

import (
	"context"
	"testing"
	"time"

	"github.com/hatena/mackerelstatsd/parser"
)

func TestServerListenAndServe(t *testing.T) {
	var s Server
	s.Error = func(err error) { t.Logf("error %v", err) }
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := s.ListenAndServe(ctx, ":0", func(addr string, m *parser.Metric) error { return nil })
	if err != nil {
		t.Fatal(err)
	}
}
