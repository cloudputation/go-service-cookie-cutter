package stats

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials"

	"github.com/cloudputation/service-seed/packages/config"
)

// tracerProvider holds the provider for graceful shutdown
var tracerProvider *sdktrace.TracerProvider

// Tracer is the global tracer for service-seed
var Tracer trace.Tracer

// InitTraces initializes the OTLP trace exporter
func InitTraces(t *config.Telemetry) error {
	ctx := context.Background()
	traces := t.Traces

	// Build exporter options
	var opts []otlptracegrpc.Option
	opts = append(opts, otlptracegrpc.WithEndpoint(traces.Endpoint))

	// TLS configuration (from top-level telemetry)
	if t.TLS != nil && t.TLS.Enabled {
		tlsCfg, err := loadTLSConfig(t.TLS)
		if err != nil {
			return fmt.Errorf("failed to load TLS config: %v", err)
		}
		opts = append(opts, otlptracegrpc.WithTLSCredentials(
			credentials.NewTLS(tlsCfg)))
	} else {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	// Headers for authentication (from top-level telemetry)
	if len(t.Headers) > 0 {
		opts = append(opts, otlptracegrpc.WithHeaders(t.Headers))
	}

	// Create exporter
	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return fmt.Errorf("failed to create trace exporter: %v", err)
	}

	// Create resource (same as metrics)
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("service-seed"),
		semconv.ServiceVersion(Version),
		semconv.DeploymentEnvironment(Environment),
	)

	// Create sampler based on config
	var sampler sdktrace.Sampler
	if traces.SamplingRate <= 0 || traces.SamplingRate >= 1.0 {
		sampler = sdktrace.AlwaysSample()
	} else {
		sampler = sdktrace.TraceIDRatioBased(traces.SamplingRate)
	}

	// Create tracer provider
	tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sampler),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tracerProvider)
	Tracer = tracerProvider.Tracer("service-seed")

	return nil
}

// ShutdownTraces gracefully shuts down the tracer provider
func ShutdownTraces(ctx context.Context) error {
	if tracerProvider == nil {
		return nil
	}
	return tracerProvider.Shutdown(ctx)
}

// Helper types for span attributes (re-exported for convenience)
var (
	StringAttribute = attribute.String
	IntAttribute    = attribute.Int
	BoolAttribute   = attribute.Bool
)

// StatusError is a convenience for setting error status on spans
var StatusError = codes.Error
