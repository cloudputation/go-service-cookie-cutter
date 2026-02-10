# stats

## Purpose
Implements application observability using OpenTelemetry with dual export support: Prometheus (pull-based) and OTLP gRPC (push-based). Provides basic instrumentation for HTTP requests and application errors. Designed to be extended with additional metrics, middleware, and distributed tracing as needed.

## Key Files
- `stats.go` (53 lines): Metrics initialization, counter definitions, dual exporter setup (Prometheus + OTLP gRPC when configured)

## Main Exports

**Initialization:**
- `InitMetrics() error`: Initializes OpenTelemetry metrics with Prometheus exporter (OTLP optional via config)

**Legacy Metrics:**
- `ErrorCounter api.Int64Counter`: Count application errors
- `HealthEndpointCounter api.Int64Counter`: Count health endpoint hits
- `SystemMetricsEndpointCounter api.Int64Counter`: Count metrics endpoint hits

## Interactions
- Depends on `config` package for telemetry configuration (OTLP endpoint, service name, environment, TLS, headers)
- Consumed by `api` package for HTTP endpoint instrumentation
- Exports metrics via Prometheus format at `/v1/system/metrics` endpoint (pull)
- Exports metrics via OTLP gRPC to configured collector (push, optional)

## Configuration/Dependencies
- Uses OpenTelemetry SDK (`go.opentelemetry.io/otel/metric`, `go.opentelemetry.io/otel/sdk/metric`)
- Uses Prometheus exporter (`go.opentelemetry.io/otel/exporters/prometheus`)
- Uses OTLP gRPC exporters (`go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc`) when configured
- Meter name: `CFS.Metrics`
- Resource attributes: `service.name`, `service.version`, `deployment.environment` (when OTLP configured)
- OTLP export enabled when `config.AppConfig.Telemetry` is configured with endpoint
- Supports mutual TLS for OTLP gRPC when certificate paths provided in config
- Supports custom headers for authentication (e.g., API keys)

## Example Usage
```go
import "github.com/organization/service-seed/packages/stats"

// Initialize metrics on startup
if err := stats.InitMetrics(); err != nil {
    log.Fatal(err)
}

// Manual metric recording
stats.ErrorCounter.Add(ctx, 1)
stats.HealthEndpointCounter.Add(ctx, 1)
```

## Implementation Notes
- All counters are Int64 (monotonically increasing)
- Counters initialized during `InitMetrics()` - must be called before use
- Thread-safe by design (OpenTelemetry handles concurrency)
- Dual export: Prometheus (pull) always enabled, OTLP gRPC (push) optional via config
- Graceful shutdown via `Shutdown()` flushes pending metrics before exit (when implemented)
- Supports mutual TLS for OTLP gRPC with client certificate, key, and CA certificate paths

## Future Enhancements
This package provides a foundation for observability. Consider adding:
- **HTTP Middleware**: Automatic request instrumentation (see sentinel/stats/middleware.go)
- **Helper Functions**: Consistent metric recording with labels (see sentinel/stats/helpers.go)
- **Distributed Tracing**: OTLP trace export and span utilities (see sentinel/stats/traces.go)
- **Additional Metrics**: Histograms for latency, gauges for resource usage, custom business metrics
- **Graceful Shutdown**: Implement `Shutdown(ctx)` to flush metrics before exit

---
Handles metrics and counters. Provides basic observability for application health and performance. Extend with middleware and tracing for comprehensive telemetry.
