# Go Application Boilerplate Guide

## Purpose
Production-ready Go template with configuration, logging, metrics, CLI, and HTTP API. **Extend, don't replace.**

## Quick Start
```bash
export SS_CONFIG_FILE_PATH=./config.hcl
./service-seed agent
```
**Endpoints**: `/v1/health`, `/v1/system/metrics` (port 8080)

## Architecture
```
main.go → config.LoadConfiguration() → stats.InitMetrics() → logger.InitLogger() → cli.SetupRootCommand()
```

## Core Packages

### Config (`packages/config/`)
HCL-based configuration with env var support.
```go
type Configuration struct {
    LogDir  string `hcl:"log_dir"`
    DataDir string `hcl:"data_dir"`
    Server  Server `hcl:"server,block"`
}
```
**Usage**: `config.LoadConfiguration()`, `config.AppConfig.Server.ServerPort`
**Env**: `SS_CONFIG_FILE_PATH` (default: `/etc/service-seed/config.hcl`)

### Logger (`packages/logger/`)
Multi-output (stdout + file) with levels: Debug, Info, Warn, Error, Fatal.
```go
l.Info("Server starting on port %s", port)
l.Error("Failed: %v", err)
```

### Stats (`packages/stats/`)
OpenTelemetry/Prometheus metrics.
```go
stats.ErrorCounter.Add(context.Background(), 1)
```
**Counters**: `ErrorCounter`, `HealthEndpointCounter`, `SystemMetricsEndpointCounter`

### CLI (`packages/cli/`)
Cobra-based. **Command**: `agent` (starts HTTP server)
```go
// Add commands in SetupRootCommand()
var newCmd = &cobra.Command{
    Use: "cmd", Short: "Description",
    Run: func(cmd *cobra.Command, args []string) { /* logic */ },
}
rootCmd.AddCommand(newCmd)
```

### API (`packages/api/`)
HTTP server with health and metrics endpoints.
```go
// Add endpoints in StartServer()
http.HandleFunc("/v1/your-endpoint", v1.YourHandler)
```

### Bootstrap (`packages/bootstrap/`)
Creates data directories from config.

## Extension Patterns

### Add Endpoint
1. Create handler in `packages/api/v1/your_feature.go`
2. Check method, increment metrics, log operation
3. Register in `packages/api/server.go`: `http.HandleFunc("/v1/path", v1.Handler)`

### Add Metrics
1. Declare counter in `packages/stats/stats.go`: `var YourCounter api.Int64Counter`
2. Initialize in `InitMetrics()`: `YourCounter, err = meter.Int64Counter("name", ...)`
3. Use: `stats.YourCounter.Add(r.Context(), 1)`

### Extend Config
```go
type Configuration struct {
    // Existing fields...
    YourFeature YourConfig `hcl:"your_feature,block"`
}
```
Update `config.hcl` accordingly.

## Rules for AI Agents

### DO
✅ Use existing packages (config, logger, stats, api, cli)
✅ Follow patterns: error handling, logging, metrics
✅ Extend existing modules
✅ See `.factory/examples/simple-todo-api.md` for complete example

### DON'T
❌ Modify core boilerplate files
❌ Duplicate functionality
❌ Skip error handling, logging, or metrics

## Development
```bash
make build         # Build binary
make docker-build  # Build container
```

**Complete Example**: See `.factory/examples/simple-todo-api.md`
**Guidelines**: `.factory/meta/SOFTWARE-ENGINEERING-GUIDELINES.md`
