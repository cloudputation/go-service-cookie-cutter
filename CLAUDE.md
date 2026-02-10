# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go application boilerplate with modular package architecture, HCL configuration, OpenTelemetry observability (metrics, logs, traces), and multi-registry Docker deployment. Designed as a template for creating new Go services with production-ready patterns.

## Build & Development Commands

```bash
# Build binary (outputs to build/iterator)
make build

# Build Docker image
make docker-build

# Full pipeline: build → docker-build → docker-push
make all

# Local development in Docker
make local-deploy              # Full rebuild and run
make local-restart             # Quick container restart (skip Go build)

# Clean build artifacts
make clean

# Show all available targets and configuration
make help
```

### Registry Configuration

The Makefile supports multiple container registries via environment variables:

```bash
# GCP Artifact Registry (default)
make docker-push

# Docker Hub
make REGISTRY=docker.io IMAGE_ORG=myorg docker-push

# AWS ECR
make REGISTRY=123456789.dkr.ecr.region.amazonaws.com IMAGE_ORG=myrepo docker-push
```

Version is read from `API_VERSION` file (fallback: 0.0.2).

## Application Initialization Order

The application follows a strict initialization sequence (see `main.go`):

1. **Load Configuration** - Parse HCL config file
2. **Initialize OTLP Logs** - Set up log export (if enabled)
3. **Initialize Logger** - Dual output (stdout/file + optional OTLP)
4. **Initialize Metrics** - Prometheus + optional OTLP export
5. **Initialize Traces** - Distributed tracing (if enabled)
6. **Start CLI** - Execute command-line interface

Graceful shutdown happens in reverse order with context timeouts.

## Configuration System

### HCL Configuration
Config file: `config.hcl` (default: `/etc/service-seed/config.hcl`)

Override via environment variable:
```bash
export SS_CONFIG_FILE_PATH=/path/to/config.hcl
```

### Configuration Structure
```hcl
log_dir = "logs"
data_dir = "data"

server {
  port = "8080"
  address = "0.0.0.0"
}

# Optional: OpenTelemetry export
telemetry {
  endpoint = "localhost:4317"  # Shared OTLP gRPC endpoint

  metrics {
    enabled = true
    protocol = "grpc"          # "grpc" or "http"
    interval_seconds = 60
  }

  logs {
    enabled = true
  }

  traces {
    enabled = true
    sampling_rate = 1.0        # 0.0-1.0
  }
}
```

### Telemetry Endpoint Inheritance
- Metrics, logs, and traces inherit the shared `telemetry.endpoint` if not explicitly overridden
- Each signal can specify its own endpoint to use different collectors
- All OTLP export uses gRPC protocol (not HTTP)

## Package Architecture

The codebase uses a **modular package structure** where each package is self-contained:

```
packages/
├── api/          - HTTP server, REST endpoints, health checks
├── bootstrap/    - Filesystem initialization (directories, state)
├── cli/          - Cobra CLI command structure
├── config/       - HCL configuration parsing (split: config.go, telemetry.go)
├── logger/       - Dual-output logging (hclog + OTLP adapter)
└── stats/        - OpenTelemetry metrics, middleware, traces
```

### Import Path Pattern
All packages use the module path:
```go
import "github.com/organization/service-seed/packages/<package>"
```

When customizing for a new service, update `go.mod` module path and all imports.

## CLAUDELET Documentation Pattern

**CRITICAL: Read package-specific documentation before modifying code**

Each package has a `CLAUDELET.md` file containing:
- Purpose and responsibilities
- Key files and exports
- Implementation patterns
- Interactions with other packages
- Configuration requirements
- Example usage

**Before modifying any package**, read its CLAUDELET.md:
```bash
cat packages/<package-name>/CLAUDELET.md
```

Available CLAUDELETs:
- `packages/api/CLAUDELET.md` - HTTP server architecture
- `packages/bootstrap/CLAUDELET.md` - Filesystem initialization patterns
- `packages/cli/CLAUDELET.md` - CLI command structure
- `packages/config/CLAUDELET.md` - Configuration loading and defaults
- `packages/logger/CLAUDELET.md` - Dual-logger pattern, OTLP export
- `packages/stats/CLAUDELET.md` - Metrics, middleware, tracing

This hierarchical documentation scales better than monolithic docs as the codebase grows.

## Observability & Telemetry

### Logging
- Uses HashiCorp `hclog` with dual output:
  - Human-readable: stdout + file (`logs/service-seed.log`)
  - JSON format: OTLP gRPC export (optional)
- Structured logging with key-value pairs
- Log levels: debug, info, warn, error, fatal
- OTLP adapter in `packages/logger/otlp_adapter.go`

### Metrics
- OpenTelemetry SDK with Prometheus exporter (always enabled)
- Optional OTLP gRPC push export when telemetry configured
- HTTP middleware for automatic request instrumentation
- Custom metrics via `stats.RecordHTTPRequest()`, `stats.RecordError()`

### Traces
- OpenTelemetry distributed tracing
- HTTP middleware automatically creates spans
- Trace context propagated through operations
- Configurable sampling rate (0.0-1.0)

### Middleware Pattern
Wrap HTTP handlers with metrics/tracing middleware:
```go
import "github.com/organization/service-seed/packages/stats"

handler := stats.MetricsMiddleware(http.HandlerFunc(myHandler))
```

## Key Design Patterns

### Configuration
- **Modular defaults**: Each config section (server, telemetry) has its own `applyDefaults()` function
- **Environment override**: `SS_CONFIG_FILE_PATH` environment variable
- **Validation**: Type checking and required field validation during parsing
- **Separation**: Core config in `config.go`, telemetry in `telemetry.go`

### Logging
- **Dual-logger pattern**: Human-readable + JSON for OTLP
- **Interface abstraction**: `Logger` interface for testing and dependency injection
- **Named loggers**: `logger.NewLogger("api")` for subsystem-specific logging
- **Deferred cleanup**: Always defer `logger.CloseLogger()` after initialization

### Observability
- **Opt-in OTLP**: Works without telemetry config; OTLP export enabled when configured
- **Graceful shutdown**: Context-based timeouts for flushing telemetry on exit
- **Resource attributes**: Service name, version, environment automatically added
- **Middleware-first**: HTTP instrumentation via middleware, not manual recording

## Module Structure Notes

### Binary Name vs Application Name
- **Binary**: `iterator` (historical name, used in `make build`)
- **Application**: `service-seed` (logical name, used in configs and logs)
- When customizing: Update `APP_NAME` in Makefile and module paths

### Versioning
- Version stored in `API_VERSION` file at repository root
- Docker tags use sanitized version (+ replaced with -)
- Version injected into telemetry resource attributes

### Build Configuration
- Target: Linux AMD64 (CGO disabled for static linking)
- Build directory: `./build/`
- Docker platform: `linux/amd64`

## Docker & Deployment

### Dockerfile
- Alpine-based (minimal attack surface)
- Non-root user (UID/GID 991)
- dumb-init for proper signal handling
- Supports ARG overrides for build-time configuration

### GKE Deployment
Optional GKE deployment via `helpers/gke-deploy.sh`:
```bash
make gke-deploy              # Deploy all environments
make gke-deploy-local        # Deploy local environment
make gke-deploy-dev          # Deploy dev environment
make gke-deploy-prod         # Deploy prod environment
```

Falls back to local Docker build if `helpers/gke-deploy.sh` not present.

## Customizing for New Services

1. **Update module path**: Change `go.mod` and all import paths from `github.com/organization/service-seed` to your service
2. **Update application name**: Change `APP_NAME` in Makefile, "SERVICE-SEED" in logger, "service-seed" in config paths
3. **Update config path**: Change `/etc/service-seed/` to `/etc/your-service/` in `config/config.go`
4. **Add business logic**: Implement handlers in `packages/api/`, add CLI commands in `packages/cli/`
5. **Update CLAUDELETs**: Modify package CLAUDELET.md files to reflect your service's purpose
6. **Configure telemetry**: Update `config.hcl` with OTLP endpoints for your observability stack

## Code Quality Standards

This codebase follows the software reliability rules defined in `~/.claude/CLAUDE.md`:
- Simple control flow (max 3 levels nesting, cyclomatic complexity < 10)
- Bounded operations (explicit loop termination, timeouts on external calls)
- Explicit resource management (defer cleanup immediately after acquisition)
- Small functions (20-30 lines, single responsibility)
- Defensive validation (fail fast, meaningful errors)
- Minimal scope (local > instance > global)
- Explicit error handling (never ignore errors)
- Minimal abstraction (direct code over indirection)
- Clear data flow (pure functions preferred, explicit parameters)
- Zero compiler/linter warnings

Refer to the global CLAUDE.md for detailed rationale on each rule.
