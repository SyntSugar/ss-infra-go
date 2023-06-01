package hooks

import (
	pro "github.com/SyntSugar/ss-infra-go/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

type performanceMetrics struct {
	Latencies *prometheus.HistogramVec
	QPS       *prometheus.CounterVec
}

var metrics *performanceMetrics

const (
	namespace = "infra"
	subsystem = "redis"
)

func setupMetrics() {
	labels := []string{"node", "command", "status"}
	buckets := prometheus.ExponentialBuckets(1, 2, 16)
	newHistogram := func(name string, labels ...string) *prometheus.HistogramVec {
		return pro.NewHistogramHelper(namespace, subsystem, name, buckets, labels...)
	}
	newCounter := func(name string, labels ...string) *prometheus.CounterVec {
		return pro.NewCounterHelper(namespace, subsystem, name, labels...)
	}
	metrics = &performanceMetrics{
		Latencies: newHistogram("latency", labels...),
		QPS:       newCounter("qps", labels...),
	}
}

func Init() {
	setupMetrics()
}
