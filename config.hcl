# Service Seed Configuration File
# This file defines the core configuration for the service-seed application

# Directory where log files will be stored
log_dir = "logs"

# Directory where application data will be stored
data_dir = "data"

# HTTP Server configuration block
server {
  # Port on which the HTTP server will listen
  port = "8080"

  # Address to bind the server (0.0.0.0 for all interfaces, 127.0.0.1 for localhost only)
  address = "0.0.0.0"
}

# Telemetry configuration (optional)
# Uncomment to enable OpenTelemetry export via OTLP gRPC
# telemetry {
#   # Shared OTLP endpoint (inherited by metrics/logs/traces if not overridden)
#   endpoint = "localhost:4317"
#
#   # Optional: Shared TLS configuration
#   # tls {
#   #   enabled = true
#   #   insecure = false
#   #   ca_file = "/path/to/ca.crt"
#   # }
#
#   # Optional: Shared headers (e.g., for authentication)
#   # headers = {
#   #   "Authorization" = "Bearer token"
#   # }
#
#   # Metrics export configuration
#   metrics {
#     enabled = true
#     protocol = "grpc"           # "grpc" or "http"
#     interval_seconds = 60       # Export interval
#     # endpoint = "localhost:4317"  # Optional: Override shared endpoint
#   }
#
#   # Logs export configuration
#   # logs {
#   #   enabled = true
#   #   # endpoint = "localhost:4317"  # Optional: Override shared endpoint
#   # }
#
#   # Traces export configuration
#   # traces {
#   #   enabled = true
#   #   sampling_rate = 1.0        # 0.0-1.0 (1.0 = 100%)
#   #   # endpoint = "localhost:4317"  # Optional: Override shared endpoint
#   # }
# }
