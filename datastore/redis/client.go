package redis

import (
	"runtime"
	"time"

	"github.com/SyntSugar/ss-infra-go/datastore"
	"github.com/SyntSugar/ss-infra-go/datastore/redis/hooks"

	"github.com/go-redis/redis/extra/redisotel/v8"
	"github.com/go-redis/redis/v8"
)

type Options struct {
	*redis.Options

	EnabledMetric    bool
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
		EnabledMetric:    true,
		EnabledOtelTrace: false,
	}
}

func NewClientWithTLS(options *Options, clientTLSConfig *datastore.ClientTLSConfig) (*redis.Client, error) {
	if options == nil {
		options = DefaultOptions()
	}
	if clientTLSConfig != nil {
		tlsConfig, err := clientTLSConfig.Build()
		if err != nil {
			return nil, err
		}
		options.TLSConfig = tlsConfig
	}
	return NewClient(options), nil
}

func NewClient(options *Options) *redis.Client {
	if options == nil {
		options = DefaultOptions()
	}
	client := redis.NewClient(options.Options)
	if options.EnabledMetric {
		client.AddHook(hooks.NewMetricsHook(client))
	}
	if options.EnabledOtelTrace {
		client.AddHook(redisotel.NewTracingHook())
	}
	return client
}
