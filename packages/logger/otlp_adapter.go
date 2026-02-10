package logger

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"google.golang.org/grpc/credentials"
)

// Build-time variables (shared with stats package)
var (
	Version     = "dev"
	Environment = "development"
)

// loggerProvider holds the OTEL logger provider for graceful shutdown
var loggerProvider *sdklog.LoggerProvider

// otlpWriter holds reference to the writer for diagnostics
var otlpWriter *otlpLogWriter

// OTLPLogsOptions configures the OTLP log exporter
type OTLPLogsOptions struct {
	// Endpoint is the OTLP collector address (e.g., "localhost:4317")
	Endpoint string
	// TLS configuration (optional)
	TLS *OTLPLogsTLSOptions
	// Headers for authentication (e.g., API keys)
	Headers map[string]string
}

// OTLPLogsTLSOptions holds TLS settings for OTLP log export
type OTLPLogsTLSOptions struct {
	Enabled  bool
	Insecure bool
	CAFile   string
	CertFile string
	KeyFile  string
}

// otlpLogWriter implements io.Writer to adapt hclog JSON output to OTLP
type otlpLogWriter struct {
	logger   log.Logger
	emitCount int64 // tracks number of log records emitted to OTLP
}

// Write implements io.Writer, parsing hclog JSON and emitting OTEL log records
func (w *otlpLogWriter) Write(p []byte) (n int, err error) {
	// hclog outputs one JSON line per log entry
	lines := strings.Split(string(p), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse the JSON log entry
		var entry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			// Not JSON, skip (shouldn't happen with hclog JSONFormat)
			continue
		}

		// Build OTEL log record
		record := log.Record{}

		// Extract and set timestamp
		if ts, ok := entry["@timestamp"].(string); ok {
			if t, err := time.Parse(time.RFC3339Nano, ts); err == nil {
				record.SetTimestamp(t)
			}
			delete(entry, "@timestamp")
		} else {
			record.SetTimestamp(time.Now())
		}

		// Extract and set severity
		if level, ok := entry["@level"].(string); ok {
			record.SetSeverity(mapHCLogLevel(level))
			record.SetSeverityText(strings.ToUpper(level))
			delete(entry, "@level")
		}

		// Extract and set message body
		if msg, ok := entry["@message"].(string); ok {
			record.SetBody(log.StringValue(msg))
			delete(entry, "@message")
		}

		// Remove @module and @caller - they become attributes
		module := ""
		if m, ok := entry["@module"].(string); ok {
			module = m
			delete(entry, "@module")
		}
		caller := ""
		if c, ok := entry["@caller"].(string); ok {
			caller = c
			delete(entry, "@caller")
		}

		// Add remaining fields as attributes
		attrs := make([]log.KeyValue, 0, len(entry)+2)
		if module != "" {
			attrs = append(attrs, log.String("logger.name", module))
		}
		if caller != "" {
			attrs = append(attrs, log.String("code.filepath", caller))
		}
		for k, v := range entry {
			attrs = append(attrs, log.String(k, fmt.Sprintf("%v", v)))
		}
		record.AddAttributes(attrs...)

		// Emit the log record
		w.logger.Emit(context.Background(), record)
		atomic.AddInt64(&w.emitCount, 1)
	}

	return len(p), nil
}

// EmitCount returns the number of log records emitted to OTLP
func (w *otlpLogWriter) EmitCount() int64 {
	return atomic.LoadInt64(&w.emitCount)
}

// mapHCLogLevel converts hclog level strings to OTEL severity
func mapHCLogLevel(level string) log.Severity {
	switch strings.ToLower(level) {
	case "trace":
		return log.SeverityTrace
	case "debug":
		return log.SeverityDebug
	case "info":
		return log.SeverityInfo
	case "warn", "warning":
		return log.SeverityWarn
	case "error":
		return log.SeverityError
	case "fatal":
		return log.SeverityFatal
	default:
		return log.SeverityInfo
	}
}

// InitOTLPLogs initializes the OTLP log exporter and returns an io.Writer
// that can be added to hclog's MultiWriter
func InitOTLPLogs(opts *OTLPLogsOptions) (*otlpLogWriter, error) {
	ctx := context.Background()

	// Build exporter options
	var exporterOpts []otlploggrpc.Option
	exporterOpts = append(exporterOpts, otlploggrpc.WithEndpoint(opts.Endpoint))

	// TLS configuration
	if opts.TLS != nil && opts.TLS.Enabled {
		tlsCfg, err := loadLogTLSConfig(opts.TLS)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS config: %v", err)
		}
		exporterOpts = append(exporterOpts, otlploggrpc.WithTLSCredentials(
			credentials.NewTLS(tlsCfg)))
	} else {
		exporterOpts = append(exporterOpts, otlploggrpc.WithInsecure())
	}

	// Headers for authentication
	if len(opts.Headers) > 0 {
		exporterOpts = append(exporterOpts, otlploggrpc.WithHeaders(opts.Headers))
	}

	// Create exporter
	exporter, err := otlploggrpc.New(ctx, exporterOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP log exporter: %v", err)
	}

	// Create resource with service attributes
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("service-seed"),
		semconv.ServiceVersion(Version),
		semconv.DeploymentEnvironment(Environment),
	)

	// Create logger provider with batch processor
	loggerProvider = sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
	)

	// Get a logger from the provider
	otelLogger := loggerProvider.Logger("service-seed")

	fmt.Fprintf(os.Stderr, "[OTLP] Log exporter initialized: endpoint=%s\n", opts.Endpoint)

	otlpWriter = &otlpLogWriter{logger: otelLogger}
	return otlpWriter, nil
}

// loadLogTLSConfig creates a TLS configuration from OTLP TLS settings
func loadLogTLSConfig(cfg *OTLPLogsTLSOptions) (*tls.Config, error) {
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

// ShutdownOTLPLogs gracefully shuts down the logger provider
func ShutdownOTLPLogs(ctx context.Context) error {
	if loggerProvider == nil {
		return nil
	}
	if otlpWriter != nil {
		fmt.Fprintf(os.Stderr, "[OTLP] Shutting down log exporter: %d log records emitted\n", otlpWriter.EmitCount())
	}
	return loggerProvider.Shutdown(ctx)
}
