package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/mackerelio/mackerel-client-go"

	"github.com/mackerelio-labs/mackerel-statsd/metricname"
	"github.com/mackerelio-labs/mackerel-statsd/parser"
	"github.com/mackerelio-labs/mackerel-statsd/statsd"
)

// Sink is a destination to collect- and to flush-metrics.
type Sink struct {
	// key is metric name.
	counterValues map[string]float64
	counterLock   sync.RWMutex

	// key is metric name.
	timerValues map[string][]float64
	timerLock   sync.RWMutex
}

func NewSink() *Sink {
	return &Sink{
		counterValues: make(map[string]float64),
		timerValues:   make(map[string][]float64),
	}
}

var sink = NewSink()

func (s *Sink) AppendValue(m *parser.Metric) {
	switch m.Type {
	case parser.TypeCounter:
		s.counterLock.Lock()
		s.counterValues[m.Name] += m.Value / m.SamplingRate
		s.counterLock.Unlock()
	case parser.TypeTimer:
		s.timerLock.Lock()
		v := m.Value / 1000 // milliseconds to seconds
		s.timerValues[m.Name] = append(s.timerValues[m.Name], v/m.SamplingRate)
		s.timerLock.Unlock()
	}
}

func (s *Sink) FlushCounterValue() []*mackerel.MetricValue {
	s.counterLock.Lock()
	metrics := make([]*mackerel.MetricValue, 0, len(s.counterValues))
	for k, v := range s.counterValues {
		metrics = append(metrics, &mackerel.MetricValue{
			Name:  k,
			Value: v,
			Time:  time.Now().Unix(),
		})
		s.counterValues[k] = 0
	}
	s.counterLock.Unlock()
	return metrics
}

func CalculateStatistics(values []float64) (min, max, average, p50, p95, p99 float64) {
	if len(values) == 0 {
		return
	}
	max = 0.0
	min = math.Inf(1)
	sum := 0.0
	for _, v := range values {
		max = math.Max(max, v)
		min = math.Min(min, v)
		sum += v
	}
	average = sum / float64(len(values))

	// TODO: implement percentile support
	return
}

func (s *Sink) FlushTimerValue() []*mackerel.MetricValue {
	now := time.Now().Unix()
	s.timerLock.Lock()
	metrics := make([]*mackerel.MetricValue, 0, len(s.timerValues)*6)
	for name, values := range s.timerValues {
		min, max, average, _, _, _ := CalculateStatistics(values)
		metrics = append(metrics, &mackerel.MetricValue{
			Name:  metricname.Join(name, "min"),
			Value: min,
			Time:  now,
		})
		metrics = append(metrics, &mackerel.MetricValue{
			Name:  metricname.Join(name, "max"),
			Value: max,
			Time:  now,
		})
		metrics = append(metrics, &mackerel.MetricValue{
			Name:  metricname.Join(name, "average"),
			Value: average,
			Time:  now,
		})
		s.timerValues[name] = nil
	}
	s.timerLock.Unlock()
	return metrics
}

type Exporter struct {
	metricNames map[string]struct{}
}

type MackerelClient interface {
	PostHostMetricValuesByHostID(hostID string, values []*mackerel.MetricValue) error
	CreateGraphDefs(params []*mackerel.GraphDefsParam) error
}

func (e *Exporter) ExportMetrics(s *Sink, hostID string, client MackerelClient) error {
	if e.metricNames == nil {
		e.metricNames = make(map[string]struct{})
	}
	counterMetrics := s.FlushCounterValue()
	timerMetrics := s.FlushTimerValue()

	var metrics []*mackerel.MetricValue
	metrics = append(metrics, counterMetrics...)
	metrics = append(metrics, timerMetrics...)

	for _, m := range timerMetrics {
		name := metricname.Group(m.Name)
		if _, ok := e.metricNames[name]; !ok {
			log.Printf("creating graphDefs: %s\n", name)
			err := client.CreateGraphDefs([]*mackerel.GraphDefsParam{
				{
					Name: name,
					Unit: "seconds",
					Metrics: []*mackerel.GraphDefsMetric{
						{Name: metricname.Join(name, "*"), DisplayName: "%1"},
					},
				},
			})
			if err != nil {
				log.Printf("creating graphdDefs: %s failed: %v\n", name, err)
				continue
			}
		}
		e.metricNames[name] = struct{}{}
	}

	log.Printf("post %d metrics\n", len(metrics))
	if len(metrics) > 0 {
		err := client.PostHostMetricValuesByHostID(hostID, metrics)
		if err != nil {
			return err
		}
	}
	return nil
}

var (
	flagInterval = flag.Duration("interval", 60*time.Second, "flush interval")
	flagHostID   = flag.String("host", "", "Post host metric values to <`hostID`>")
	flagAddr     = flag.String("addr", ":8125", "`addr`ess to listen")
)

func main() {
	flag.Parse()
	hostID := *flagHostID
	if hostID == "" {
		log.Fatalln("must specify -host")
	}
	apiKey := os.Getenv("MACKEREL_APIKEY")
	if apiKey == "" {
		log.Fatalln("must set MACKEREL_APIKEY")
	}
	client := mackerel.NewClient(apiKey)

	ticker := time.NewTicker(*flagInterval)
	defer ticker.Stop()
	var e Exporter
	go func() {
		for {
			select {
			case <-ticker.C:
				err := e.ExportMetrics(sink, hostID, client)
				if err != nil {
					// TODO: retry
					log.Println(err)
				}
			}
		}
	}()

	var server statsd.Server
	server.Error = func(err error) {
		log.Println(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	defer stop()
	server.ListenAndServe(ctx, *flagAddr, func(addr string, m *parser.Metric) error {
		sink.AppendValue(m)
		fmt.Printf("From: %v Receiving metric: %+v\n", addr, m.Value)
		return nil
	})
}
