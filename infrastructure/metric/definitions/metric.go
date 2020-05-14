package definitions

import "time"

// MetricInterface as a contract
type MetricInterface interface {
	Count(name string, value int64, tags []string, rate float64) error
	Gauge(name string, value float64, tags []string, rate float64) error
	Histogram(name string, startTime time.Time, tags []string) error
}
