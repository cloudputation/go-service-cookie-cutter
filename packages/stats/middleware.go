package stats

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// MetricsMiddleware wraps an HTTP handler with automatic metrics and tracing
func MetricsMiddleware(endpoint string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Start trace span (if tracing enabled)
		var span trace.Span
		ctx := r.Context()
		if Tracer != nil {
			ctx, span = Tracer.Start(ctx, fmt.Sprintf("HTTP %s %s", r.Method, endpoint))
			defer span.End()
			r = r.WithContext(ctx)
		}

		// Wrap ResponseWriter to capture status code
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default to 200
		}

		// Call the actual handler
		next(wrapped, r)

		// Record metrics after handler completes
		duration := time.Since(start)
		RecordHTTPRequest(ctx, r.Method, endpoint, wrapped.statusCode, duration)

		// Set span attributes after handler completes
		if span != nil {
			span.SetAttributes(
				attribute.Int("http.status_code", wrapped.statusCode),
				attribute.String("http.method", r.Method),
				attribute.String("http.route", endpoint),
			)
			if wrapped.statusCode >= 400 {
				span.SetStatus(codes.Error, http.StatusText(wrapped.statusCode))
			}
		}
	}
}

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code before writing
func (w *responseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// Write ensures status code is captured even if WriteHeader isn't called explicitly
func (w *responseWriter) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

// Flush implements http.Flusher if the underlying ResponseWriter supports it
func (w *responseWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Hijack implements http.Hijacker for WebSocket upgrades
func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}
