package tracing

import (
	"context"

	"github.com/SyntSugar/ss-infra-go/consts"
)

func WithAmTraceID(parent context.Context, traceID string) context.Context {
	if parent == nil {
		return nil
	}
	return context.WithValue(parent, consts.ContextKeyTraceID, traceID)
}

func GetAmTraceID(ctx context.Context) string {
	traceID, _ := ctx.Value(consts.ContextKeyTraceID).(string)
	return traceID
}

func WithCloudflareRayID(parent context.Context, cloudflareRayID string) context.Context {
	if parent == nil {
		return nil
	}

	return context.WithValue(parent, consts.ContextKeyCloudflareRay, cloudflareRayID)
}

func GetCloudflareRayID(ctx context.Context) string {
	cloudflareRayID, _ := ctx.Value(consts.ContextKeyCloudflareRay).(string)
	return cloudflareRayID
}
