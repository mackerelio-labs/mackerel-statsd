package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/hatena/mackerelstatsd/metricname"
	"github.com/hatena/mackerelstatsd/statsd"
)

func metricName(path string) string {
	name := strings.ReplaceAll(path, "/", "_")
	return name
}

func HTTPHandler(next http.Handler, c *statsd.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		delta := time.Since(start)
		name := metricname.Join("http.request" + metricName(r.URL.Path))
		c.Timer(name, delta, &statsd.Options{SamplingRate: 1.0})
	})
}
