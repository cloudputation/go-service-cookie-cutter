package stats

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"google.golang.org/grpc/credentials"

	"github.com/cloudputation/service-seed/packages/config"
)

const meterName = "service-seed"

// Build-time variables injected via ldflags
var (
	Version     = "dev"
	Environment = "development"
)

// Meter is exported for use by helper functions
var Meter api.Meter

// meterProvider holds the provider for graceful shutdown
var meterProvider *metric.MeterProvider

// ============================================================================
// LEGACY METRICS (kept for backward compatibility)
// ============================================================================

var (
	// Deprecated: Use ErrorsTotal with labels instead
	ErrorCounter                 api.Int64Counter
	HealthEndpointCounter        api.Int64Counter
	SystemMetricsEndpointCounter api.Int64Counter
)

// ============================================================================
// HTTP METRICS
// ============================================================================

var (
	// HTTPRequestsTotal counts requests with method, endpoint, status_code labels
	HTTPRequestsTotal api.Int64Counter

	// HTTPRequestDuration measures request latency with method, endpoint labels
	HTTPRequestDuration api.Float64Histogram
)

// ============================================================================
// ERROR METRICS
// ============================================================================

var (
	// ErrorsTotal counts errors with component and error_type labels
	ErrorsTotal api.Int64Counter
)

// ============================================================================
// INITIALIZATION
// ============================================================================

func InitMetrics() error {
	// Create resource with service attributes
	// Note: Using NewWithAttributes directly to avoid schema version conflicts
	// between resource.Default() (v1.37.0) and semconv (v1.24.0)
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("service-seed"),
		semconv.ServiceVersion(Version),
		semconv.DeploymentEnvironment(Environment),
	)

	// Always create Prometheus exporter (backward compatibility)
	prometheusExporter, err := prometheus.New()
	if err != nil {
		return fmt.Errorf("failed to initialize prometheus exporter: %v", err)
	}

	// Build provider options
	providerOpts := []metric.Option{
		metric.WithResource(res),
		metric.WithReader(prometheusExporter),
	}

	// Conditionally add OTLP exporter if enabled
	if config.AppConfig.Telemetry != nil &&
		config.AppConfig.Telemetry.Metrics != nil &&
		config.AppConfig.Telemetry.Metrics.Enabled {

		otlpReader, err := createOTLPReader(config.AppConfig.Telemetry)
		if err != nil {
			return fmt.Errorf("failed to create OTLP exporter: %v", err)
		}
		providerOpts = append(providerOpts, metric.WithReader(otlpReader))
	}

	// Create meter provider with all readers
	meterProvider = metric.NewMeterProvider(providerOpts...)
	Meter = meterProvider.Meter(meterName)

	// Start Go runtime metrics collection (goroutines, memory, GC)
	if err := runtime.Start(runtime.WithMeterProvider(meterProvider)); err != nil {
		return fmt.Errorf("failed to start runtime metrics: %v", err)
	}

	if err := initLegacyMetrics(); err != nil {
		return err
	}
	if err := initHTTPMetrics(); err != nil {
		return err
	}
	if err := initErrorMetrics(); err != nil {
		return err
	}
	if err := initGaugeMetrics(); err != nil {
		return err
	}

	return nil
}

func initLegacyMetrics() error {
	var err error

	ErrorCounter, err = Meter.Int64Counter(
		"agent_errors",
		api.WithDescription("Counts the number of errors during agent runtime (deprecated: use service_errors_total)"),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize error counter: %v", err)
	}

	HealthEndpointCounter, err = Meter.Int64Counter(
		"health_endpoint_hits",
		api.WithDescription("Counts the number of hits to the /health endpoint"),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize health endpoint counter: %v", err)
	}

	SystemMetricsEndpointCounter, err = Meter.Int64Counter(
		"system_metrics_endpoint_hits",
		api.WithDescription("Counts the number of hits to the /system/metrics endpoint"),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize system metrics endpoint counter: %v", err)
	}

	return nil
}

func initHTTPMetrics() error {
	var err error

	HTTPRequestsTotal, err = Meter.Int64Counter(
		"service_http_requests_total",
		api.WithDescription("HTTP requests by method, endpoint, and status code"),
		api.WithUnit("{request}"),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize service_http_requests_total: %v", err)
	}

	HTTPRequestDuration, err = Meter.Float64Histogram(
		"service_http_request_duration_seconds",
		api.WithDescription("HTTP request latency"),
		api.WithUnit("s"),
		api.WithExplicitBucketBoundaries(0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize service_http_request_duration_seconds: %v", err)
	}

	return nil
}

func initErrorMetrics() error {
	var err error

	ErrorsTotal, err = Meter.Int64Counter(
		"service_errors_total",
		api.WithDescription("Errors by component and type"),
		api.WithUnit("{error}"),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize service_errors_total: %v", err)
	}

	return nil
}

func initGaugeMetrics() error {
	// Placeholder for future gauge metrics
	return nil
}

// ============================================================================
// OTLP EXPORTER
// ============================================================================

// createOTLPReader creates a PeriodicReader with OTLP gRPC exporter
func createOTLPReader(t *config.Telemetry) (metric.Reader, error) {
	ctx := context.Background()
	metrics := t.Metrics

	var opts []otlpmetricgrpc.Option
	opts = append(opts, otlpmetricgrpc.WithEndpoint(metrics.Endpoint))

	// TLS configuration (from top-level telemetry)
	if t.TLS != nil && t.TLS.Enabled {
		tlsCfg, err := loadTLSConfig(t.TLS)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS config: %v", err)
		}
		opts = append(opts, otlpmetricgrpc.WithTLSCredentials(
			credentials.NewTLS(tlsCfg)))
	} else {
		opts = append(opts, otlpmetricgrpc.WithInsecure())
	}

	// Headers for authentication (from top-level telemetry)
	if len(t.Headers) > 0 {
		opts = append(opts, otlpmetricgrpc.WithHeaders(t.Headers))
	}

	exporter, err := otlpmetricgrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %v", err)
	}

	// PeriodicReader handles background export goroutine
	interval := time.Duration(metrics.IntervalSeconds) * time.Second
	return metric.NewPeriodicReader(exporter, metric.WithInterval(interval)), nil
}

// loadTLSConfig creates a TLS configuration from OTLP TLS settings
func loadTLSConfig(cfg *config.OTLPTLSConfig) (*tls.Config, error) {
	tlsCfg := &tls.Config{
		InsecureSkipVerify: cfg.Insecure,
	}

	// Load client certificate if provided (mutual TLS)
	if cfg.CertFile != "" && cfg.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load client cert/key: %v", err)
		}
		tlsCfg.Certificates = []tls.Certificate{cert}
	}

	// Load CA certificate if provided
	if cfg.CAFile != "" {
		caCert, err := os.ReadFile(cfg.CAFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA file: %v", err)
		}
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}
		tlsCfg.RootCAs = caCertPool
	}

	return tlsCfg, nil
}

// Shutdown gracefully shuts down the meter provider, flushing any pending metrics
func Shutdown(ctx context.Context) error {
	if meterProvider == nil {
		return nil
	}
	return meterProvider.Shutdown(ctx)
}
