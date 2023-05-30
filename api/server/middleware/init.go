package middleware

import (
	prome "github.com/SyntSugar/ss-infra-go/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

type serverMetrics struct {
	Latencies        *prometheus.HistogramVec
	HTTPCodes        *prometheus.CounterVec
	Payload          *prometheus.CounterVec
	HTTPServerPanics *prometheus.CounterVec
}

var serMetrics *serverMetrics

const (
	namespace = "infra"
	subsystem = "http_api"
)

func setupMetrics() {
	labels := []string{"host", "uri", "method", "code", "custom"}
	buckets := prometheus.ExponentialBuckets(1, 2, 16)
	newHistogram := func(name string, labels ...string) *prometheus.HistogramVec {
		return prome.NewHistogramHelper(namespace, subsystem, name, buckets, labels...)
	}
	newCounter := func(name string, labels ...string) *prometheus.CounterVec {
		return prome.NewCounterHelper(namespace, subsystem, name, labels...)
	}
	serMetrics = &serverMetrics{
		Latencies:        newHistogram("request_latency", labels...),
		HTTPCodes:        newCounter("http_code", labels...),
		Payload:          newCounter("http_payload", labels...),
		HTTPServerPanics: newCounter("http_server_panic"),
	}
}

func init() {
	setupMetrics()
}
