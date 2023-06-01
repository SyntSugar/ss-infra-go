package hooks

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/SyntSugar/ss-infra-go/consts"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

type MetricsHook struct {
	client *redis.Client
}

func NewMetricsHook(client *redis.Client) *MetricsHook {
	return &MetricsHook{client: client}
}

func (hook *MetricsHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		ctx = context.WithValue(ctx, consts.ContextStartTimeKey, time.Now())
		next(ctx, cmd)
		hook.record(ctx, cmd.Name(), cmd.Err())
		return nil
	}
}

func (hook *MetricsHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		ctx = context.WithValue(ctx, consts.ContextStartTimeKey, time.Now())
		next(ctx, cmds)
		var firstErr error
		for _, cmd := range cmds {
			if cmd.Err() != nil {
				firstErr = cmd.Err()
				break
			}
		}
		hook.record(ctx, "pipeline", firstErr)
		return nil
	}
}

func (hook *MetricsHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(_ context.Context, network, addr string) (net.Conn, error) {
		conn, err := net.Dial(network, addr)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
}

func (hook MetricsHook) record(ctx context.Context, cmd string, err error) {
	startTime, ok := ctx.Value(consts.ContextStartTimeKey).(time.Time)
	if !ok {
		return
	}
	durationMS := time.Since(startTime).Milliseconds()
	status := "ok"
	if err != nil && !errors.Is(err, redis.Nil) {
		status = "error"
	}
	labels := prometheus.Labels{"node": hook.client.Options().Addr, "command": cmd, "status": status}
	metrics.QPS.With(labels).Inc()
	metrics.Latencies.With(labels).Observe(float64(durationMS))
}
