package statsd

import "time"

type Client struct {
	backend Driver
}

type Driver interface {
	Inc(name string, n int64, rate float64, tags []*Tag) error
	Timer(name string, d time.Duration, rate float64, tags []*Tag) error
	Close() error
}

type Tag struct {
	Key   string
	Value string
}

type Options struct {
	SamplingRate float64
	Tags         []*Tag
}

func NewClient(d Driver) *Client {
	return &Client{backend: d}
}

func (c *Client) Inc(name string, n int64, opts *Options) error {
	return c.backend.Inc(name, n, opts.SamplingRate, opts.Tags)
}

func (c *Client) Timer(name string, d time.Duration, opts *Options) error {
	return c.backend.Timer(name, d, opts.SamplingRate, opts.Tags)
}

func (c *Client) Close() error {
	return c.backend.Close()
}
