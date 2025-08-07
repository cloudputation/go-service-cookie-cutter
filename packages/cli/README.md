# cli

## Purpose
Defines the application's command-line interface (CLI), providing commands for config validation, agent startup, and system status checks.

## Main Exports
- `SetupRootCommand() *cobra.Command`: Returns the root Cobra command with subcommands.
- CLI commands: `config`, `agent`, `system status`.

## Interactions
- Uses `bootstrap` to initialize the agent/filesystem.
- Uses `api` to start the HTTP server.
- Uses `logger` and `stats` for logging and metrics.
- Uses `config` for config validation.

## Configuration/Dependencies
- Uses Cobra for CLI parsing.
- Relies on other local modules for command implementations.

## Example Usage
```
go
import "github.com/organization/service-seed/packages/cli"
rootCmd := cli.SetupRootCommand()
rootCmd.Execute()
```

---
Focuses on CLI orchestration and command definitions.
