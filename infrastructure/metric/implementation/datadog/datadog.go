package datadog

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

// Datadog to hold datadog client state
type Datadog struct {
	client *statsd.Client
}

// New init new datadog client
func New(serviceName, env, source string) *Datadog {
	datadog, err := statsd.New(source)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Get hostname
	host, err := os.Hostname()
	if err != nil {
		host = "undefined"
	}

	// Get service name
	if len(serviceName) < 1 {
		log.Fatal(errors.New("Datadog service name should be provided"))
	}

	datadog.Namespace = fmt.Sprintf("enterprise_%s.", serviceName)
	datadog.Tags = append(datadog.Tags, "env:"+env, "host:"+host)

	log.Println("Datadog initialized...")

	return &Datadog{
		client: datadog,
	}
}

// Count tracks how many times something happened per second
func (datadog *Datadog) Count(name string, value int64, tags []string, rate float64) error {
	err := datadog.client.Count(name, value, tags, rate)
	if err != nil {
		return err
	}
	return nil
}

// Gauge measures the value of a metric at a particular time
func (datadog *Datadog) Gauge(name string, value float64, tags []string, rate float64) error {
	err := datadog.client.Gauge(name, value, tags, rate)
	if err != nil {
		return err
	}
	return nil
}

// Histogram tracks the statistical distribution of a set of values on each host
func (datadog *Datadog) Histogram(name string, startTime time.Time, tags []string) error {
	elapsedTime := time.Since(startTime).Seconds() * 1000
	err := datadog.client.Histogram(name, elapsedTime, tags, float64(1))
	if err != nil {
		return err
	}
	return nil
}
