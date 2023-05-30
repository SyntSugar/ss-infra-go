package middleware

import (
	"fmt"

	"github.com/SyntSugar/ss-infra-go/consts"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const tracerKey = "otel-tracer"

// NewOpenTelemetryTracing returns a Gin middleware function for tracing incoming requests.
// If no propagator or tracerProvider is provided, it uses the global ones.
func NewOpenTelemetryTracing(serviceName string, propagator propagation.TextMapPropagator, tracerProvider oteltrace.TracerProvider) gin.HandlerFunc {
	if propagator == nil {
		propagator = otel.GetTextMapPropagator()
	}
	if tracerProvider == nil {
		tracerProvider = otel.GetTracerProvider()
	}
	tracer := tracerProvider.Tracer(consts.OtelDefaultTracerName)

	// The returned function is the actual Gin middleware.
	return func(c *gin.Context) {
		c.Set(tracerKey, tracer)
		savedCtx := c.Request.Context()
		defer func() {
			c.Request = c.Request.WithContext(savedCtx)
		}()
		ctx := propagator.Extract(savedCtx, propagation.HeaderCarrier(c.Request.Header))
		opts := []oteltrace.SpanStartOption{
			oteltrace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", c.Request)...),
			oteltrace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(c.Request)...),
			oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(serviceName, c.FullPath(), c.Request)...),
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		}
		path := c.FullPath()
		if path == "" {
			path = fmt.Sprintf("HTTP %s route not found", c.Request.Method)
		}
		spanName := fmt.Sprintf("%s %s", c.Request.Method, path)
		sCtx, span := tracer.Start(ctx, spanName, opts...)
		defer span.End()

		// pass the span through the request context
		c.Request = c.Request.WithContext(sCtx)

		// serve the request to the next middleware
		c.Next()

		status := c.Writer.Status()
		attrs := semconv.HTTPAttributesFromHTTPStatusCode(status)
		spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(status)
		span.SetAttributes(attrs...)
		span.SetStatus(spanStatus, spanMessage)
		if len(c.Errors) > 0 {
			span.SetAttributes(attribute.String("gin.errors", c.Errors.String()))
		}
	}
}
