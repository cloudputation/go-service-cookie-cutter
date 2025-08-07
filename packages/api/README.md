# api

## Purpose
Provides HTTP API endpoints for the application, exposing health, system status, and metrics routes for external monitoring and orchestration.

## Main Exports
- `StartServer()`: Starts the HTTP server, registers endpoints `/v1/health`, `/v1/system/status`, `/v1/system/metrics`.
- `v1/health.go`: `HealthHandler` (GET /v1/health)
- `v1/system_status.go`: `SystemStatusHandler`, `SystemStatusHandlerWrapper` (GET /v1/system/status)

## Interactions
- Uses `config` for configuration (e.g., server port).
- Uses `logger` for logging.
- Uses `stats` for metrics and counters.
- Uses `consul` for distributed state (system status).

## Configuration/Dependencies
- Relies on configuration from the `config` package (server port, etc.).
- Expects Prometheus client (`promhttp`).
- Reads API version from `./VERSION` file for health endpoint.

## Example Usage
```
go
import "github.com/organization/service-seed/packages/api"
api.StartServer()
```

---
This module is focused on API logic only. See sibling modules for CLI, configuration, logging, and metrics.
