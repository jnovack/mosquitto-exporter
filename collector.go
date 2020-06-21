package main

import (
	"sync"
	"time"

	version "github.com/jnovack/go-version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var (
	ignoreKeyMetrics = map[string]string{
		"$SYS/broker/timestamp":        "The timestamp at which this particular build of the broker was made. Static.",
		"$SYS/broker/version":          "The version of the broker. Static.",
		"$SYS/broker/clients/active":   "//deprecated// in favour of $SYS/broker/clients/connected",
		"$SYS/broker/clients/inactive": "//deprecated// in favour of $SYS/broker/clients/disconnected",
	}
	counterKeyMetrics = map[string]string{
		"$SYS/broker/bytes/received":            "The total number of bytes received since the broker started.",
		"$SYS/broker/bytes/sent":                "The total number of bytes sent since the broker started.",
		"$SYS/broker/messages/received":         "The total number of messages of any type received since the broker started.",
		"$SYS/broker/messages/sent":             "The total number of messages of any type sent since the broker started.",
		"$SYS/broker/publish/bytes/received":    "The total number of PUBLISH bytes received since the broker started.",
		"$SYS/broker/publish/bytes/sent":        "The total number of PUBLISH bytes sent since the broker started.",
		"$SYS/broker/publish/messages/received": "The total number of PUBLISH messages received since the broker started.",
		"$SYS/broker/publish/messages/sent":     "The total number of PUBLISH messages sent since the broker started.",
		"$SYS/broker/publish/messages/dropped":  "The total number of PUBLISH messages that have been dropped due to inflight/queuing limits.",
		"$SYS/broker/uptime":                    "The total number of seconds since the broker started.",
		"$SYS/broker/clients/maximum":           "The maximum number of clients connected simultaneously since the broker started",
		"$SYS/broker/clients/total":             "The total number of clients connected since the broker started.",
	}
	counterGaugeMetrics = map[string]string{
		"$SYS/broker/clients/connected":    "number of clients connected",
		"$SYS/broker/clients/disconnected": "number of clients disconnected",
		"$SYS/broker/subscriptions/count":  "number of active subscriptions",
	}
)

func timeTrack(ch chan<- prometheus.Metric, start time.Time, name string) {
	elapsed := time.Since(start)
	log.Debugf("%s took %.3fs", name, float64(elapsed.Milliseconds())/1000)

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc("go_task_time", "Go task elasped time", []string{}, prometheus.Labels{"task": name, "application": version.Application}),
		prometheus.GaugeValue,
		float64(elapsed.Milliseconds())/1000,
	)
}

// Collector TODO Comment
type Collector struct {
	desc string
}

// Metric TODO Comment
type Metric struct {
	name   string
	help   string
	value  float64
	labels map[string]string
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {

	metrics := make(chan prometheus.Metric)
	go func() {
		c.Collect(metrics)
		close(metrics)
	}()
	for m := range metrics {
		ch <- m.Desc()
	}
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {

	wg := sync.WaitGroup{}

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(parseForPrometheus(version.Application), "github.com/jnovack/"+version.Application, []string{}, prometheus.Labels{"version": version.Version}),
		prometheus.GaugeValue,
		1,
	)

	// Datacenter Metrics
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer timeTrack(ch, time.Now(), "datacenterMetrics")
		cm := datacenterMetrics(ch)
		for _, m := range cm {

			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(m.name, m.help, []string{}, m.labels),
				prometheus.GaugeValue,
				float64(m.value),
			)
		}

	}()

	wg.Wait()
}

// NewCollector TODO Comment
func NewCollector() *Collector {
	return &Collector{
		desc: "mosquitto Collector",
	}
}

func datacenterMetrics(ch chan<- prometheus.Metric) []Metric {
	// defer cancel()

	var metrics []Metric
	metrics = append(metrics, Metric{name: "vmware_datastore_size", help: "Maximum capacity of this datastore, in bytes.", value: float64(1), labels: map[string]string{"datastore": "a", "cluster": "a", "datacenter": "a"}})

	return metrics
}
