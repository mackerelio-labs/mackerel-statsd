package main

import (
	"io"
	"log"
	"net/http"

	"github.com/mackerelio-labs/mackerel-statsd/example/driver/cactusstatsd"
	"github.com/mackerelio-labs/mackerel-statsd/example/middleware"
	"github.com/mackerelio-labs/mackerel-statsd/statsd"
)

func helloHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello, world!\n")
}

func main() {
	d, err := cactusstatsd.NewDriver("localhost:8125", "custom.statsd.sample")
	if err != nil {
		log.Fatal(err)
	}
	c := statsd.NewClient(d)
	defer c.Close()
	http.Handle("/", middleware.HTTPHandler(http.HandlerFunc(helloHandler), c))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
