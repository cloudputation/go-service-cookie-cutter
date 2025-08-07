# stats

## Purpose
Implements application metrics and statistics tracking, exposing Prometheus counters and generating system state.

## Main Exports
- `InitMetrics() error`: Initializes Prometheus metrics.
- `FileSystemInfo` struct: Describes filesystem state.
- `GenerateState() error`: Generates and stores system state summary.
- `GetFileSystemInfo(dirPath string) (string, error)`: Reads and summarizes filesystem info.
- Various Prometheus counters (e.g., `ErrorCounter`, `HealthEndpointCounter`).

## Interactions
- Uses `config` for directory paths.
- Uses `consul` for distributed state storage.

## Configuration/Dependencies
- Uses OpenTelemetry and Prometheus libraries.
- Expects directories as defined in config.

## Example Usage
```
go
import "github.com/organization/service-seed/packages/stats"
stats.InitMetrics()
```

---
Handles metrics, counters, and system state generation only.
