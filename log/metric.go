package log

import (
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap/zapcore"
)

type metric struct {
	samplingCounter *prometheus.CounterVec
}

// logSamplingMetrics is a hook function that logs the metrics associated with
// the sampling decision of a log entry. It increases a Prometheus counter labelled
// with the log entry's level and the sampling decision ("sampled" or "dropped").
func logSamplingMetrics(entry zapcore.Entry, decision zapcore.SamplingDecision) {
	var decisionLabel string
	if decision == zapcore.LogDropped {
		decisionLabel = "dropped"
	} else {
		decisionLabel = "sampled"
	}

	metricsLabels := prometheus.Labels{
		"level":    entry.Level.String(),
		"decision": decisionLabel,
	}

	globalMetric.samplingCounter.With(metricsLabels).Inc()
}
