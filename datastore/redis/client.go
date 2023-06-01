package redis

import (
	"context"
	"runtime"
	"time"

	"github.com/SyntSugar/ss-infra-go/datastore/redis/hooks"
	"github.com/uptrace/uptrace-go/uptrace"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type Options struct {
	*redis.Options

	EnabledOtelMetric    bool
	EnabledOtelTrace bool
}

func DefaultOptions() *Options {
	return &Options{
		Options: &redis.Options{
			DialTimeout:  1200 * time.Millisecond,
			ReadTimeout:  1500 * time.Millisecond,
			WriteTimeout: time.Second,
			MinIdleConns: runtime.NumCPU(),
		},
		EnabledOtelMetric: true,
		EnabledOtelTrace:  false,
	}
}



func NewClient(options *Options) *redis.Client {
	if options == nil {
		options = DefaultOptions()
	}

	rdb := redis.NewClient(options.Options)

	hooks.Init() 

	if options.EnabledOtelMetric {
		metricsHook := hooks.NewMetricsHook(rdb)
		rdb.AddHook(metricsHook)
	}
	
	if options.EnabledOtelTrace {
		uptrace.ConfigureOpentelemetry(
			// copy your project DSN here or use UPTRACE_DSN env var
			uptrace.WithDSN("http://project2_secret_token@localhost:14317/2"),
	
			uptrace.WithServiceName("myservice"),
			uptrace.WithServiceVersion("v1.0.0"),
		)
		defer uptrace.Shutdown(context.Background())

		if err := redisotel.InstrumentTracing(rdb); err != nil {
			panic(err)
		}
	}

	return rdb
}
