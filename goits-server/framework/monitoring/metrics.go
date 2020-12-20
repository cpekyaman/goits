package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var registry *prometheus.Registry
var factory promauto.Factory

func init() {
	registry = prometheus.NewRegistry()
	factory = promauto.With(registry)
}

// HistogramBundle represents a group of histograms for the same metric with labels.
type HistogramBundle struct {
	h *prometheus.HistogramVec
}

// Histogram is a single labeled metric instance of a bundle.
type Histogram struct {
	o prometheus.Observer
}

// With returns a Histogram for the given labels.
func (this HistogramBundle) With(labels map[string]string) Histogram {
	return Histogram{this.h.With(labels)}
}

// Record observes the given value as part of metric collection.
func (this Histogram) Record(value float64) {
	this.o.Observe(value)
}

// NewHistogram creates a new histogram metric with the given name and expecting provided labels.
func NewHistogram(name string, labels []string) HistogramBundle {
	return HistogramBundle{
		h: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    name,
			Buckets: prometheus.LinearBuckets(25, 25, 10),
		}, labels),
	}
}

// CounterBundle represents a group of counters for the same metric with labels.
type CounterBundle struct {
	c *prometheus.CounterVec
}

// Counter is a single labeled metric instance of a bundle.
type Counter struct {
	c prometheus.Counter
}

// With returns a Counter for the given labels.
func (this CounterBundle) With(labels map[string]string) Counter {
	return Counter{this.c.With(labels)}
}

// Incr increments the counter value by one.
func (this Counter) Incr() {
	this.c.Inc()
}

// Add adds the provided value to counter.
func (this Counter) Add(value float64) {
	this.c.Add(value)
}

// NewCounter creates a new counter metric with the given name and expecting provided labels.
func NewCounter(name string, labels []string) CounterBundle {
	return CounterBundle{
		c: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: name,
		}, labels),
	}
}
