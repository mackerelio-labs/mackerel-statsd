package cactusstatsd

import (
	"time"

	"github.com/cactus/go-statsd-client/v5/statsd"

	proto "github.com/mackerelio-labs/mackerel-statsd/statsd"
)

type Driver struct {
	c statsd.Statter
}

var _ proto.Driver = (*Driver)(nil)

func NewDriver(addr, prefix string) (*Driver, error) {
	c, err := statsd.NewClientWithConfig(&statsd.ClientConfig{
		Address: addr,
		Prefix:  prefix,
	})
	if err != nil {
		return nil, err
	}
	return &Driver{c: c}, nil
}

func (d *Driver) Inc(name string, n int64, rate float64, tags []*proto.Tag) error {
	rawTags := make([]statsd.Tag, len(tags))
	for i, tag := range tags {
		rawTags[i] = statsd.Tag{tag.Key, tag.Value}
	}
	return d.c.Inc(name, n, float32(rate), rawTags...)
}

func (d *Driver) Timer(name string, dur time.Duration, rate float64, tags []*proto.Tag) error {
	rawTags := make([]statsd.Tag, len(tags))
	for i, tag := range tags {
		rawTags[i] = statsd.Tag{tag.Key, tag.Value}
	}
	return d.c.TimingDuration(name, dur, float32(rate), rawTags...)
}

func (d *Driver) Close() error {
	return d.c.Close()
}
