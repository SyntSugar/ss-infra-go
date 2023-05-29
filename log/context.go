package log

import (
	"context"

	"github.com/SyntSugar/ss-infra-go/consts"
	"github.com/SyntSugar/ss-infra-go/tracing"

	"go.uber.org/zap"
)

func AppendContextFields(parent context.Context, fields ...zap.Field) context.Context {
	if parent == nil {
		return nil
	}

	// The value-passing method is used here, so that the subclass function sets new fields,
	// which will not affect the fields set by the caller, causing some unpredictable behavior
	oldFields, ok := parent.Value(ctxLogger).([]zap.Field)
	if ok {
		fields = append(fields, oldFields...)
	}

	//nolint:staticcheck,SA1029
	ctx := context.WithValue(parent, ctxLogger, fields)

	return ctx
}

func GetContextFields(ctx context.Context) []zap.Field {
	if ctx == nil {
		return nil
	}

	fields, ok := ctx.Value(ctxLogger).([]zap.Field)
	if !ok {
		return nil
	}

	if traceID := tracing.GetAmTraceID(ctx); traceID != "" {
		fields = append(fields, zap.String(consts.KeyAMTraceID, traceID))
	}

	if cloudflareRayID := tracing.GetCloudflareRayID(ctx); cloudflareRayID != "" {
		fields = append(fields, zap.String(consts.KeyCloudflareRay, cloudflareRayID))
	}

	return fields
}

// DynamicDebugLogging open the debug level logging in dynamically
func DynamicDebugLogging(parent context.Context) context.Context {
	if parent == nil {
		return parent
	}
	return context.WithValue(parent, consts.ContextKeyEnableDebugLogging, true)
}
