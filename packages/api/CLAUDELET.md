# API - HTTP Server & REST Endpoints

## Purpose

Provides HTTP server with REST API endpoints for health checks and system metrics. Central entry point for all client interactions with the service. Includes Prometheus metrics export and basic health monitoring.

## Core Components

### HTTP Server
- **Port**: Configured via `server.port` (default: 3001)
- **Router**: Standard `http.HandleFunc` registration
- **Middleware**: CORS, logging (when implemented)
- **Graceful shutdown**: Context-based termination with signal handling (when implemented)

### Endpoint Registration

All endpoints registered in `server.go` StartServer():

**Health & Metrics**:
- `GET /v1/health` - Health check endpoint
- `GET /v1/system/metrics` - Prometheus metrics

## Key Files

**Server Initialization**:
- `server.go` (29 lines) - HTTP server setup, endpoint registration

**v1/ Package** (API v1):
- `health.go` - Health check HTTP handler

## Exports

**Main Server**:
- `StartServer()` - Initialize and start HTTP server

**v1 Exports**:
- `HealthHandler()` - Health check endpoint

**HTTP Handlers**:
- `HealthHandler()` - Health check endpoint

## Dependencies

- **config** - Server configuration (port, address)
- **logger** - HTTP request logging with structured key-value pairs
- **stats** - Metrics tracking (endpoint counters, Prometheus)

## Initialization Flow

```go
StartServer():
  1. Load configuration
  2. Register HTTP endpoints
  3. Start HTTP server
```

## Configuration

```hcl
server {
  port = "3001"
  address = "0.0.0.0"
}
```

## Thread Safety

- **HTTP Handlers**: Stateless, thread-safe

## Error Handling

- Invalid requests: 400 Bad Request (when validation implemented)
- Not found: 404 Not Found
- Internal errors: 500 Internal Server Error

## Metrics

Basic metrics exported via `/v1/system/metrics`:
- `health_endpoint_hits` - Health endpoint hits
- `system_metrics_endpoint_hits` - Metrics endpoint hits
- `agent_errors` - Application errors

When telemetry is configured, metrics are exported to both Prometheus (scrape endpoint) and OTLP gRPC collector.

See [packages/stats/CLAUDELET.md](../stats/CLAUDELET.md) for metrics details.

## Design Patterns

**Simple HTTP Server**:
- Standard library `net/http` for HTTP server
- Direct endpoint registration with `http.HandleFunc`
- Prometheus exporter for metrics
- Extensible for additional endpoints and middleware

## Future Enhancements

Consider adding:
- **Middleware**: HTTP metrics and tracing middleware (see sentinel/api/server.go)
- **WebSocket Support**: Real-time streaming (see sentinel/api/v1/websocket.go)
- **Graceful Shutdown**: Signal handling and context-based shutdown
- **Request Validation**: Input validation and error handling
- **CORS Support**: Cross-origin resource sharing configuration
- **Rate Limiting**: Request throttling and abuse prevention
- **Authentication**: API key or JWT-based auth
