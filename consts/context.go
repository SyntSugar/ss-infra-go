package consts

// The reason is that using a built-in type as a key can potentially cause collisions
// if another package decides to use the same type and value as a key.
// To avoid this, you can define your own type and use that as the key instead.
type ContextKey string

const (
	ContextKeyTraceID            ContextKey = "traceID"
	ContextKeyCloudflareRay      ContextKey = "cloudflareRay"
	ContextKeyRequestPath        ContextKey = "requestPath"
	ContextKeyEnableDebugLogging ContextKey = "enableDebugLog"
	ContextKeyMetricLabel        ContextKey = "metricLabel"
	ContextStartTimeKey          ContextKey = "startTime"
	ContextSegmentKey            ContextKey = "segment"
)
