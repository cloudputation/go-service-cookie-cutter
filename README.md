# Application Boilerplate

## Overview
This project provides a modular boilerplate for Go applications with a focus on service orchestration, configuration, metrics, logging, and extensible CLI/API patterns. Each core module is documented below.

## Table of Contents
- [api](./packages/api/README.md)
- [bootstrap](./packages/bootstrap/README.md)
- [cli](./packages/cli/README.md)
- [config](./packages/config/README.md)
- [logger](./packages/logger/README.md)
- [stats](./packages/stats/README.md)

---

## Module Summaries

### [api](./packages/api/README.md)
- **Purpose:** HTTP API endpoints for health, system status, and metrics.
- **Exports:** `StartServer()`, `HealthHandler`, `SystemStatusHandler`, `SystemStatusHandlerWrapper`.
- **Interactions:** Uses `config`, `logger`, `stats`, `consul`.
- **Example:**
  ```go
  import "github.com/organization/service-seed/packages/api"
  api.StartServer()
  ```

### [bootstrap](./packages/bootstrap/README.md)
- **Purpose:** Initializes and bootstraps the filesystem agent and Consul state.
- **Exports:** `BootstrapFileSystem()`
- **Interactions:** Uses `config`, `consul`, `logger`, `stats`.
- **Example:**
  ```go
  import "github.com/organization/service-seed/packages/bootstrap"
  err := bootstrap.BootstrapFileSystem()
  ```

### [cli](./packages/cli/README.md)
- **Purpose:** Command-line interface for config validation, agent startup, and system status.
- **Exports:** `SetupRootCommand()`
- **Interactions:** Uses `bootstrap`, `api`, `logger`, `stats`, `config`.
- **Example:**
  ```go
  import "github.com/organization/service-seed/packages/cli"
  rootCmd := cli.SetupRootCommand()
  rootCmd.Execute()
  ```

### [config](./packages/config/README.md)
- **Purpose:** Centralized configuration loading and access.
- **Exports:** `Configuration`, `LoadConfiguration()`, `GetConfigPath()`, `AppConfig`.
- **Interactions:** Used by all modules.
- **Example:**
  ```go
  import "github.com/organization/service-seed/packages/config"
  err := config.LoadConfiguration()
  ```

### [logger](./packages/logger/README.md)
- **Purpose:** Centralized logging utilities for file/stdout and levels.
- **Exports:** `InitLogger()`, `CloseLogger()`, `Debug`, `Info`, `Warn`, `Error`, `Fatal`.
- **Interactions:** Used by all modules.
- **Example:**
  ```go
  import "github.com/organization/service-seed/packages/logger"
  logger.InitLogger("/var/log/service-seed", "info")
  logger.Info("Started")
  ```

### [stats](./packages/stats/README.md)
- **Purpose:** Metrics and statistics tracking, Prometheus counters, system state.
- **Exports:** `InitMetrics()`, `FileSystemInfo`, `GenerateState()`, `GetFileSystemInfo()`
- **Interactions:** Uses `config`, `consul`.
- **Example:**
  ```go
  import "github.com/organization/service-seed/packages/stats"
  stats.InitMetrics()
  ```

---

For detailed module documentation, see each module's README.
