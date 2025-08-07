# logger

## Purpose
Provides centralized logging utilities for the application, supporting file and stdout output with log levels.

## Main Exports
- `InitLogger(logDirPath, logLevelController string) error`: Initializes logger.
- `CloseLogger()`: Closes log file.
- `Debug`, `Info`, `Warn`, `Error`, `Fatal`: Logging helpers.

## Interactions
- Used by all modules for logging.

## Configuration/Dependencies
- Uses HashiCorp hclog for logging.
- Log level controlled by `logLevelController` argument.

## Example Usage
```
go
import "github.com/organization/service-seed/packages/logger"
logger.InitLogger("/var/log/service-seed", "info")
logger.Info("Started")
```

---
Handles only logging logic for the application.
