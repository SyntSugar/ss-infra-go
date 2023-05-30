package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/SyntSugar/ss-infra-go/consts"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func CollectMetrics(c *gin.Context) {
	startTime := time.Now()
	c.Next()
	latency := time.Since(startTime).Milliseconds()

	uri := c.FullPath()
	// uri was empty means not found routes, so rewrite it to /not_found here
	if c.Writer.Status() == http.StatusNotFound && uri == "" {
		uri = "/not_found"
	}
	customMetricLabel := "-"
	if v, ok := c.Get(string(consts.ContextKeyMetricLabel)); ok {
		if label, ok := v.(string); ok {
			customMetricLabel = label
		}
	}
	labels := prometheus.Labels{
		"host":   c.Request.Host,
		"uri":    uri,
		"method": c.Request.Method,
		"code":   strconv.Itoa(c.Writer.Status()),
		"custom": customMetricLabel,
	}
	serMetrics.HTTPCodes.With(labels).Inc()
	serMetrics.Latencies.With(labels).Observe(float64(latency))
	size := c.Writer.Size()
	if size > 0 {
		serMetrics.Payload.With(labels).Add(float64(size))
	}
}
