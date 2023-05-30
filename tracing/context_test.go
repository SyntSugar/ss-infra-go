package tracing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAmTraceID(t *testing.T) {
	ctx := context.Background()

	// Test traceID insertion and retrieval
	traceID := "trace123"
	ctx = WithAmTraceID(ctx, traceID)
	assert.Equal(t, traceID, GetAmTraceID(ctx))

	// Test traceID replacement
	newTraceID := "trace456"
	ctx = WithAmTraceID(ctx, newTraceID)
	assert.Equal(t, newTraceID, GetAmTraceID(ctx))

	// Test that nil context returns nil
	var nilCtx context.Context
	assert.Nil(t, WithAmTraceID(nilCtx, traceID))
}

func TestCloudflareRayID(t *testing.T) {
	ctx := context.Background()

	// Test cloudflareRayID insertion and retrieval
	cloudflareRayID := "cloudflareRay123"
	ctx = WithCloudflareRayID(ctx, cloudflareRayID)
	assert.Equal(t, cloudflareRayID, GetCloudflareRayID(ctx))

	// Test cloudflareRayID replacement
	newCloudflareRayID := "cloudflareRay456"
	ctx = WithCloudflareRayID(ctx, newCloudflareRayID)
	assert.Equal(t, newCloudflareRayID, GetCloudflareRayID(ctx))

	// Test that nil context returns nil
	var nilCtx context.Context
	assert.Nil(t, WithCloudflareRayID(nilCtx, cloudflareRayID))
}
