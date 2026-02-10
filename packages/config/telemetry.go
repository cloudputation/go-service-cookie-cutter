package config

// Telemetry holds telemetry export configuration
type Telemetry struct {
	// Shared config (inherited by metrics, logs, traces)
	Endpoint string            `hcl:"endpoint,optional"`
	TLS      *OTLPTLSConfig    `hcl:"tls,block"`
	Headers  map[string]string `hcl:"headers,optional"`

	// Signal-specific config
	Metrics *OTLPMetricsConfig `hcl:"metrics,block"`
	Logs    *OTLPLogsConfig    `hcl:"logs,block"`
	Traces  *OTLPTracesConfig  `hcl:"traces,block"`
}

// OTLPMetricsConfig holds OTLP metrics exporter configuration
type OTLPMetricsConfig struct {
	// Enabled controls whether OTLP metric export is active (default: false)
	Enabled bool `hcl:"enabled"`

	// Endpoint overrides the shared telemetry endpoint for metrics
	Endpoint string `hcl:"endpoint,optional"`

	// Protocol specifies the transport protocol: "grpc" (default) or "http"
	Protocol string `hcl:"protocol,optional"`

	// IntervalSeconds is the export interval in seconds (default: 60)
	IntervalSeconds int `hcl:"interval_seconds,optional"`
}

// OTLPLogsConfig holds OTLP logs exporter configuration
type OTLPLogsConfig struct {
	// Enabled controls whether OTLP log export is active (default: false)
	Enabled bool `hcl:"enabled"`

	// Endpoint overrides the shared telemetry endpoint for logs
	Endpoint string `hcl:"endpoint,optional"`
}

// OTLPTracesConfig holds OTLP traces exporter configuration
type OTLPTracesConfig struct {
	// Enabled controls whether OTLP trace export is active (default: false)
	Enabled bool `hcl:"enabled"`

	// Endpoint overrides the shared telemetry endpoint for traces
	Endpoint string `hcl:"endpoint,optional"`

	// SamplingRate controls trace sampling (0.0-1.0, where 1.0 = 100%)
	// Default: 1.0 (sample all traces)
	SamplingRate float64 `hcl:"sampling_rate,optional"`
}

// OTLPTLSConfig holds TLS settings for OTLP export
type OTLPTLSConfig struct {
	// Enabled enables TLS for the connection
	Enabled bool `hcl:"enabled"`

	// Insecure skips certificate verification (not recommended for production)
	Insecure bool `hcl:"insecure,optional"`

	// CAFile is the path to CA certificate for server verification
	CAFile string `hcl:"ca_file,optional"`

	// CertFile is the path to client certificate (for mutual TLS)
	CertFile string `hcl:"cert_file,optional"`

	// KeyFile is the path to client key (for mutual TLS)
	KeyFile string `hcl:"key_file,optional"`
}

// applyTelemetryDefaults sets default values for telemetry configuration
func applyTelemetryDefaults() {
	if AppConfig.Telemetry == nil {
		return
	}

	t := AppConfig.Telemetry

	// Apply metrics defaults
	if t.Metrics != nil {
		// Inherit endpoint from top-level if not specified
		if t.Metrics.Endpoint == "" {
			t.Metrics.Endpoint = t.Endpoint
		}

		// Default protocol is gRPC
		if t.Metrics.Protocol == "" {
			t.Metrics.Protocol = "grpc"
		}

		// Default export interval is 60 seconds
		if t.Metrics.IntervalSeconds == 0 {
			t.Metrics.IntervalSeconds = 60
		}
	}

	// Apply logs defaults
	if t.Logs != nil && t.Logs.Enabled {
		// Inherit endpoint from top-level if not specified
		if t.Logs.Endpoint == "" {
			t.Logs.Endpoint = t.Endpoint
		}
	}

	// Apply traces defaults
	if t.Traces != nil && t.Traces.Enabled {
		// Inherit endpoint from top-level if not specified
		if t.Traces.Endpoint == "" {
			t.Traces.Endpoint = t.Endpoint
		}

		// Default sampling rate is 1.0 (100%)
		if t.Traces.SamplingRate == 0 {
			t.Traces.SamplingRate = 1.0
		}
	}
}
