package consts

type ContextKey string

const (
	ContextKeyTraceID            ContextKey = "traceID"
	ContextKeyCloudflareRay      ContextKey = "cloudflareRay"
	ContextKeyRequestPath        ContextKey = "requestPath"
	ContextKeyEnableDebugLogging ContextKey = "enableDebugLog"
	ContextKeyMetricLabel        ContextKey = "metricLabel"
)
