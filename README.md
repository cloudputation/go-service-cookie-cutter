# Application Boilerplate

Go service boilerplate with production-ready patterns: modular packages, HCL configuration, OpenTelemetry observability, and multi-registry Docker deployment.

## Quick Start

```bash
# Build the application
make build

# Run locally
make local-deploy

# See all available commands
make help
```

## What's Inside

**Configuration** - HCL-based config with sensible defaults and environment variable overrides

**Logging** - Structured logging with optional OpenTelemetry export (stdout + file + OTLP)

**Metrics** - Prometheus metrics with optional OTLP push export, HTTP middleware for automatic instrumentation

**Tracing** - Distributed tracing via OpenTelemetry with configurable sampling

**CLI** - Cobra-based command-line interface with extensible command structure

**Docker** - Multi-registry support (GCP Artifact Registry, Docker Hub, AWS ECR) with security best practices

## Project Structure

```
packages/
├── api/          HTTP server and REST endpoints
├── bootstrap/    Filesystem initialization
├── cli/          Command-line interface
├── config/       HCL configuration management
├── logger/       Centralized logging with OTLP support
└── stats/        Metrics, middleware, and tracing
```

## Documentation

- **[CLAUDE.md](./CLAUDE.md)** - Architecture guide for AI assistants working on this codebase
- **Package CLAUDELETs** - Each package has a `CLAUDELET.md` with implementation details:
  - [api/CLAUDELET.md](./packages/api/CLAUDELET.md)
  - [bootstrap/CLAUDELET.md](./packages/bootstrap/CLAUDELET.md)
  - [cli/CLAUDELET.md](./packages/cli/CLAUDELET.md)
  - [config/CLAUDELET.md](./packages/config/CLAUDELET.md)
  - [logger/CLAUDELET.md](./packages/logger/CLAUDELET.md)
  - [stats/CLAUDELET.md](./packages/stats/CLAUDELET.md)

## Configuration

Edit `config.hcl` or set `SS_CONFIG_FILE_PATH` environment variable:

```hcl
log_dir = "logs"
data_dir = "data"

server {
  port = "8080"
  address = "0.0.0.0"
}

# Optional: Enable OpenTelemetry export
telemetry {
  endpoint = "localhost:4317"

  metrics { enabled = true }
  logs { enabled = true }
  traces { enabled = true }
}
```

## Customizing for Your Service

1. Update module path in `go.mod`
2. Change `APP_NAME` in `GNUmakefile`
3. Update import paths throughout codebase
4. Modify package logic for your use case
5. Update CLAUDELETs to document your changes

See [CLAUDE.md](./CLAUDE.md) for detailed customization instructions.
