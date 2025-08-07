# bootstrap

## Purpose
Initializes and bootstraps the application's underlying filesystem agent, including Consul state and initial metrics.

## Main Exports
- `BootstrapFileSystem() error`: Bootstraps filesystem state, initializes Consul, refreshes state, logs progress.

## Interactions
- Uses `config` for configuration.
- Uses `consul` for distributed state init and bootstrap.
- Uses `logger` for logging.
- Uses `stats` to generate and refresh state.

## Configuration/Dependencies
- Reads from `config.AppConfig` and `config.RootDir`.
- Requires Consul and stats modules.

## Example Usage
```
go
import "github.com/organization/service-seed/packages/bootstrap"
err := bootstrap.BootstrapFileSystem()
```

---
Handles only the initial agent/filesystem bootstrapping logic.
