package tracing

import (
	"context"
	"os"
	"time"

	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"google.golang.org/grpc"
)

type protocol string

const (
	GRPC protocol = "gRPC"
	HTTP protocol = "HTTP"
)

type OTLConfig struct {
	Protocol protocol
	Endpoint string
	Sampler  sdktrace.Sampler
	Exporter string
	URLPath  string
	IsExport bool
}

var DefaultOTLConfig = &OTLConfig{
	Protocol: GRPC,
	Endpoint: getDefaultOTLEndpoint(),
	Sampler:  sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0.1)),
	Exporter: os.Getenv("OTEL_EXPORTER_NAME"),
	URLPath:  "/ote/v1/traces",
	IsExport: true,
}

// InitOTLProvider initializes an OTLP exporter, and configures the corresponding trace providers.
func InitOTLProvider(config *OTLConfig) (func() error, error) {
	if config == nil {
		config = DefaultOTLConfig
	}
	ctx := context.Background()
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(config.Exporter),
		),
	)
	if err != nil {
		return nil, err
	}

	otel.SetTextMapPropagator(b3.New())
	if !config.IsExport {
		return func() error {
			return nil
		}, nil
	}

	var traceExporter *otlptrace.Exporter
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	if config.Protocol == GRPC {
		// Set up a trace exporter
		traceExporter, err = otlptracegrpc.New(timeoutCtx,
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(config.Endpoint),
			otlptracegrpc.WithDialOption(grpc.WithBlock()),
		)
	} else {
		traceExporter, err = otlptrace.New(context.Background(),
			otlptracehttp.NewClient(
				otlptracehttp.WithEndpoint(config.Endpoint),
				otlptracehttp.WithInsecure(),
				otlptracehttp.WithEndpoint(config.URLPath),
			),
		)
	}
	if err != nil {
		return nil, err
	}

	// Register the trace exporter with a TracerProvider
	// using a batch span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(config.Sampler),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	return func() error {
		return tracerProvider.Shutdown(context.Background())
	}, err
}

// TODO: When integrate other envs endpoint with k8s.
func getDefaultOTLEndpoint() string {
	return "localhost:4317"
}
