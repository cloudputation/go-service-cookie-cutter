package stats

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	api "go.opentelemetry.io/otel/metric"
)

// ============================================================================
// HTTP HELPERS
// ============================================================================

// RecordHTTPRequest records HTTP request count and duration metrics
func RecordHTTPRequest(ctx context.Context, method, endpoint string, statusCode int, duration time.Duration) {
	attrs := api.WithAttributes(
		attribute.String("method", method),
		attribute.String("endpoint", endpoint),
		attribute.Int("status_code", statusCode),
	)
	HTTPRequestsTotal.Add(ctx, 1, attrs)

	durationAttrs := api.WithAttributes(
		attribute.String("method", method),
		attribute.String("endpoint", endpoint),
	)
	HTTPRequestDuration.Record(ctx, duration.Seconds(), durationAttrs)
}

// ============================================================================
// ERROR HELPERS
// ============================================================================

// RecordError records an error with component and type labels
func RecordError(ctx context.Context, component, errorType string) {
	ErrorsTotal.Add(ctx, 1,
		api.WithAttributes(
			attribute.String("component", component),
			attribute.String("error_type", errorType),
		),
	)
	// Also increment legacy counter for backward compatibility
	ErrorCounter.Add(ctx, 1)
}

// ============================================================================
// TIMER UTILITY
// ============================================================================

// Timer is a helper for measuring operation durations
type Timer struct {
	start time.Time
}

// NewTimer creates a new timer starting from now
func NewTimer() *Timer {
	return &Timer{start: time.Now()}
}

// Elapsed returns the duration since the timer was created
func (t *Timer) Elapsed() time.Duration {
	return time.Since(t.start)
}

// ObserveDuration records the elapsed time to a histogram
func (t *Timer) ObserveDuration(ctx context.Context, histogram api.Float64Histogram, attrs ...attribute.KeyValue) {
	histogram.Record(ctx, t.Elapsed().Seconds(), api.WithAttributes(attrs...))
}
