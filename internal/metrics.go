package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	prome "github.com/SyntSugar/ss-infra-go/prometheus"
)

type Metrics struct {
	Panic *prometheus.CounterVec
}

const (
	namespace = "infra"
	subsystem = "go"
)

var (
	initOnce sync.Once
	metrics  *Metrics
)

func setupMetrics() {
	metrics = &Metrics{
		Panic: prome.NewCounterHelper(namespace, subsystem, "panic", "type"),
	}
}

// Get would return the Metrics instance
func Get() *Metrics {
	initOnce.Do(setupMetrics)
	return metrics
}
