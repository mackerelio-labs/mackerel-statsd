package main

import (
	"io"
	"log"
	"net/http"

	"github.com/hatena/mackerelstatsd/example/driver/cactusstatsd"
	"github.com/hatena/mackerelstatsd/example/middleware"
	"github.com/hatena/mackerelstatsd/statsd"
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
