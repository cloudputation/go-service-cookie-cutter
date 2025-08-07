# config

## Purpose
Centralizes application configuration loading, parsing, and access for all modules.

## Main Exports
- `Configuration` struct: Holds all config fields (log dir, data dir, server, consul).
- `LoadConfiguration() error`: Loads and parses config from HCL file.
- `GetConfigPath() string`: Returns config file path.
- `AppConfig`, `ConfigPath`, `RootDir`: Global config variables.

## Interactions
- Used by all modules for config values.
- Reads environment variable `SS_CONFIG_FILE_PATH`.

## Configuration/Dependencies
- Uses Viper and HashiCorp HCL libraries for parsing.
- Expects config in `/etc/service-seed/config.hcl` or as set in env.

## Example Usage
```
go
import "github.com/organization/service-seed/packages/config"
err := config.LoadConfiguration()
```

---
Handles config parsing and global config state only.
