package main

import (
	"testing"

	"github.com/mackerelio/mackerel-client-go"
)

type mackerelClientMock struct {
	postHostMetricCalledNum int
}

func (c *mackerelClientMock) PostHostMetricValuesByHostID(hostID string, values []*mackerel.MetricValue) error {
	c.postHostMetricCalledNum++
	return nil
}

func (c *mackerelClientMock) CreateGraphDefs(params []*mackerel.GraphDefsParam) error {
	return nil
}

func TestExporterExportMetricsEmpty(t *testing.T) {
	var e Exporter
	s := NewSink()
	var c mackerelClientMock
	err := e.ExportMetrics(s, "xxx", &c)
	if err != nil {
		t.Fatal(err)
	}
	if c.postHostMetricCalledNum != 0 {
		t.Errorf("PostHostMetric should not be called, but it called %d times", c.postHostMetricCalledNum)
	}
}
