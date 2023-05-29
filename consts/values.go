package consts

const (
	HeaderAMTraceID          = "am-trace-id"
	HeaderXCloudTraceContext = "x-cloud-trace-context"
	HeaderCloudflareRay      = "CF-Ray"
	HeaderEnableDebugLogging = "Enable-Debug-Log"

	KeyAMTraceID          = "am_trace_id"
	KeyCloudflareRay      = "cloudflare_ray"
	KeyEnableDebugLogging = "enable_debug_log"

	OtelDefaultTracerName             = "default"
	OtelDefaultSpanNamePubSubProducer = "PubSub Publish"
	OtelDefaultSpanNamePubSubConsumer = "PubSub Consume"

	DefaultAccessLogPattern = `%{2006-01-02T15:04:05.999-0700}t "%{CF-Ray}i" "${AM-Trace-ID}" ` +
		`%a %A %{Host}i "%r" %s - %T "%{X-Real-IP}i" "%{X-Forwarded-For}i" ` +
		`%{Content-Length}i - %{Content-Length}o %b`
	JSONAccessLogPattern = `{"message": "AccessLogger %r %s [${AM-Trace-ID}]",` +
		`"@timestamp":"%{2006-01-02T15:04:05.999-0700}t",` +
		`"category":"http_access_log",` +
		`"context_cloudflare_ray":"%{CF-Ray}i",` +
		`"context_trace_id":"${AM-Trace-ID}",` +
		`"remote_addr":"%a",` +
		`"server_addr":"%A",` +
		`"host":"%{Host}i",` +
		`"request":"%r",` +
		`"status":"%s",` +
		`"first_byte_commit_time":"%F ms",` +
		`"request_time":%D,` +
		`"http_x_real_ip":"%{X-Real-IP}i",` +
		`"http_x_forwarded_for":"%{X-Forwarded-For}i",` +
		`"content_length":"${Content-Length}",` +
		`"body_bytes_sent":"%B bytes"}`
)
