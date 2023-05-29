package prome

import (
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// NewHistogramHelper was used to fast create and register prometheus histogram metric
func NewHistogramHelper(ns, subsystem, name string, buckets []float64, labels ...string) *prometheus.HistogramVec {
	ns = strings.ReplaceAll(ns, "-", "_")
	subsystem = strings.ReplaceAll(subsystem, "-", "_")
	name = strings.ReplaceAll(name, "-", "_")
	opts := prometheus.HistogramOpts{}
	opts.Name = name
	opts.Help = name
	opts.Namespace = ns
	opts.Subsystem = subsystem
	opts.Buckets = buckets
	histogram := prometheus.NewHistogramVec(opts, labels)
	prometheus.MustRegister(histogram)
	return histogram
}

// NewCounterHelper was used to fast create and register prometheus counter metric
func NewCounterHelper(ns, subsystem, name string, labels ...string) *prometheus.CounterVec {
	ns = strings.ReplaceAll(ns, "-", "_")
	subsystem = strings.ReplaceAll(subsystem, "-", "_")
	opts := prometheus.CounterOpts{}
	opts.Name = name
	opts.Help = name
	opts.Namespace = ns
	opts.Subsystem = subsystem
	counters := prometheus.NewCounterVec(opts, labels)
	prometheus.MustRegister(counters)
	return counters
}

// NewGaugeHelper was used to fast create and register prometheus gauge metric
func NewGaugeHelper(ns, subsystem, name string, labels ...string) *prometheus.GaugeVec {
	opts := prometheus.GaugeOpts{}
	opts.Name = name
	opts.Help = name
	opts.Namespace = strings.ReplaceAll(ns, "-", "_")
	opts.Subsystem = strings.ReplaceAll(subsystem, "-", "_")
	gauge := prometheus.NewGaugeVec(opts, labels)
	prometheus.MustRegister(gauge)
	return gauge
}

// NewSummaryHelper was used to fast create and register prometheus summary metric
func NewSummaryHelper(ns, subsystem, name string, ageBuckets, bufCap uint32, maxAge time.Duration,
	objs map[float64]float64, labels ...string) *prometheus.SummaryVec {
	opts := prometheus.SummaryOpts{
		Namespace:  strings.ReplaceAll(ns, "-", "_"),
		Subsystem:  strings.ReplaceAll(subsystem, "-", "_"),
		Name:       name,
		Help:       name,
		AgeBuckets: ageBuckets,
		BufCap:     bufCap,
		MaxAge:     maxAge,
		Objectives: objs,
	}
	summary := prometheus.NewSummaryVec(opts, labels)
	prometheus.MustRegister(summary)
	return summary
}

// DefaultSummaryHelper was used to fast create and register prometheus summary metric with default
func DefaultSummaryHelper(ns, subsystem, name string, labels ...string) *prometheus.SummaryVec {
	opts := prometheus.SummaryOpts{
		Name:      name,
		Help:      name,
		Namespace: strings.ReplaceAll(ns, "-", "_"),
		Subsystem: strings.ReplaceAll(subsystem, "-", "_"),
	}
	summary := prometheus.NewSummaryVec(opts, labels)
	prometheus.MustRegister(summary)
	return summary
}
