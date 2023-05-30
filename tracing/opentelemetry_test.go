package tracing

import (
	"context"
	"testing"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TraceExporter interface {
	ExportSpans(context.Context, []*sdktrace.ReadOnlySpan) error
	Shutdown(context.Context) error
}

type MockTraceExporter struct {
	mock.Mock
}

func (m *MockTraceExporter) ExportSpans(ctx context.Context, spans []*sdktrace.ReadOnlySpan) error {
	args := m.Called(ctx, spans)
	return args.Error(0)
}

func (m *MockTraceExporter) Shutdown(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestWithAmTraceIDAndGetAmTraceID(t *testing.T) {
	ctx := context.Background()
	traceID := "12345"
	newCtx := WithAmTraceID(ctx, traceID)

	assert.Equal(t, traceID, GetAmTraceID(newCtx))
}

func TestWithCloudflareRayIDAndGetCloudflareRayID(t *testing.T) {
	ctx := context.Background()
	cloudflareRayID := "67890"
	newCtx := WithCloudflareRayID(ctx, cloudflareRayID)

	assert.Equal(t, cloudflareRayID, GetCloudflareRayID(newCtx))
}
