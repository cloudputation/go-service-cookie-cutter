# GPT Agent Guide: Go Application Boilerplate

## 🎯 Purpose
This boilerplate serves as a **production-ready Go application template** with built-in service orchestration, configuration management, metrics collection, logging, and extensible CLI/API patterns. Use this guide to intelligently leverage existing components rather than recreating them.

## 📁 Project Structure Overview

```
application-boilerplate/
├── .factory/                   # Meta-configuration for project structure
│   ├── directory.hcl           # HCL definition of project layout
│   └── meta/                   # Guidelines and documentation
├── main.go                     # Application entry point
├── packages/                   # Core modular components
│   ├── api/                    # HTTP server and endpoints
│   ├── bootstrap/              # System initialization
│   ├── cli/                    # Command-line interface
│   ├── config/                 # Configuration management
│   ├── logger/                 # Centralized logging
│   └── stats/                  # Metrics and monitoring
├── helpers/                    # Development utilities
├── Dockerfile                  # Container configuration
└── GNUmakefile                 # Build automation
```

## 🏗️ Core Architecture

### Entry Point Flow
```go
main.go → config.LoadConfiguration() → stats.InitMetrics() → logger.InitLogger() → cli.SetupRootCommand()
```

### Key Components Integration
- **Config**: Centralized HCL-based configuration with environment variable support
- **Logger**: Multi-output logging (stdout + file) with configurable levels
- **Stats**: OpenTelemetry/Prometheus metrics collection
- **CLI**: Cobra-based command structure with subcommands
- **API**: HTTP server with health, status, and metrics endpoints
- **Bootstrap**: System initialization and Consul integration

## 📦 Package Deep Dive

### 🔧 Config Package (`packages/config/`)
**Purpose**: Centralized configuration management using HCL format

**Key Exports**:
- `AppConfig Configuration` - Global config instance
- `LoadConfiguration() error` - Loads config from file/env
- `ConfigPath string` - Path to config file
- `RootDir string` - Application root directory

**Configuration Structure**:
```go
type Configuration struct {
    LogDir  string `hcl:"log_dir"`
    DataDir string `hcl:"data_dir"`
    Server  Server `hcl:"server,block"`
    Consul  Consul `hcl:"consul,block"`
}
```

**Usage Pattern**:
**You need to follow the software engineering guidelines described at `meta/SOFTWARE-ENGINEERING-GUIDELINES.md`**.

```go
import "github.com/organization/service-seed/packages/config"
err := config.LoadConfiguration()
serverPort := config.AppConfig.Server.ServerPort
```

### 📝 Logger Package (`packages/logger/`)
**Purpose**: Structured logging with multiple outputs and configurable levels

**Key Exports**:
- `InitLogger(logDirPath, logLevel string) error`
- `Debug/Info/Warn/Error/Fatal(format string, v ...interface{})`
- `CloseLogger()`

**Features**:
- Multi-writer output (stdout + file)
- HashiCorp hclog integration
- Configurable log levels
- Automatic log file management

**Usage Pattern**:
```go
import l "github.com/organization/service-seed/packages/logger"
l.Info("Server starting on port %s", port)
l.Error("Failed to connect: %v", err)
```

### 📊 Stats Package (`packages/stats/`)
**Purpose**: OpenTelemetry metrics collection with Prometheus export

**Key Exports**:
- `InitMetrics() error` - Initialize metrics system
- `ErrorCounter api.Int64Counter` - Error tracking
- `HealthEndpointCounter api.Int64Counter` - Health endpoint hits
- `GenerateState() error` - Generate filesystem state

**Usage Pattern**:
```go
import "github.com/organization/service-seed/packages/stats"
stats.ErrorCounter.Add(context.Background(), 1)
```

### 🖥️ CLI Package (`packages/cli/`)
**Purpose**: Cobra-based command-line interface with extensible subcommands

**Key Structure**:
- Root command: `service-seed`
- Subcommands: `config`, `agent`, `status`
- Built-in configuration validation
- Agent lifecycle management

**Extension Pattern**:
```go
var newCmd = &cobra.Command{
    Use:   "newcommand",
    Short: "Description",
    Run: func(cmd *cobra.Command, args []string) {
        // Your logic here
    },
}
rootCmd.AddCommand(newCmd)
```

### 🌐 API Package (`packages/api/`)
**Purpose**: HTTP server with standard endpoints and middleware

**Built-in Endpoints**:
- `/v1/health` - Health check
- `/v1/system/status` - System status
- `/v1/system/metrics` - Prometheus metrics

**Extension Pattern**:
```go
http.HandleFunc("/v1/your-endpoint", YourHandler)
```

### 🚀 Bootstrap Package (`packages/bootstrap/`)
**Purpose**: System initialization and external service integration

**Key Functions**:
- Consul initialization and bootstrapping
- Filesystem state generation
- Service discovery setup

## 🎯 Development Patterns

### Adding New Features
1. **Identify the appropriate package** for your feature
2. **Extend existing modules** rather than creating new ones
3. **Follow established patterns** for error handling and logging
4. **Use existing configuration** and metrics systems

### Creating New Endpoints
```go
// In packages/api/v1/your_feature.go
func YourFeatureHandler(w http.ResponseWriter, r *http.Request) {
    stats.YourFeatureCounter.Add(context.Background(), 1)
    l.Info("Processing your feature request")
    
    // Your logic here
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// In packages/api/server.go
http.HandleFunc("/v1/your-feature", v1.YourFeatureHandler)
```

### Adding New CLI Commands
```go
// In packages/cli/cli.go
var yourCmd = &cobra.Command{
    Use:   "your-command",
    Short: "Description of your command",
    Run: func(cmd *cobra.Command, args []string) {
        l.Info("Executing your command")
        // Your logic here
    },
}
rootCmd.AddCommand(yourCmd)
```

### Configuration Extension
```go
// In packages/config/config.go
type Configuration struct {
    // Existing fields...
    YourFeature YourFeatureConfig `hcl:"your_feature,block"`
}

type YourFeatureConfig struct {
    Setting1 string `hcl:"setting1"`
    Setting2 int    `hcl:"setting2"`
}
```

## 🔍 Factory System (`.factory/`)

### Directory.hcl
Defines the complete project structure in HCL format. This serves as:
- **Documentation** of expected file layout
- **Validation** for project structure
- **Template** for new projects

### Meta Guidelines
- `boilerplate-guidelines.md` - Rules for GPT agents
- `software-engineering-guidelines.md` - Code quality standards

## 🚨 Critical Rules for GPT Agents

### DO:
✅ **Leverage existing packages** - Use config, logger, stats, etc.
✅ **Follow established patterns** - Error handling, logging conventions
✅ **Extend, don't replace** - Add to existing modules
✅ **Use proper imports** - Follow the established import patterns
✅ **Implement complete features** - Include error handling, logging, metrics

### DON'T:
❌ **Rewrite core boilerplate files** - Extend them instead
❌ **Create duplicate functionality** - Use existing logger, config, etc.
❌ **Skip error handling** - All functions must handle errors properly
❌ **Ignore logging** - Use the established logging patterns
❌ **Forget metrics** - Instrument new features with counters

## 🎯 Quick Start Checklist

When creating a new application:

1. **Understand the requirement** - What does the app need to do?
2. **Map to existing packages** - Which components can be reused?
3. **Plan the extension points** - Where will new code be added?
4. **Follow the patterns** - Use established error handling, logging, metrics
5. **Test integration** - Ensure new code works with existing systems
6. **Update documentation** - Add README entries for new features

## 📚 Package Import Patterns

```go
// Standard imports
import (
    "context"
    "fmt"
    "net/http"
)

// Third-party imports
import (
    "github.com/spf13/cobra"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// Local package imports
import (
    "github.com/organization/service-seed/packages/config"
    l "github.com/organization/service-seed/packages/logger"  // Aliased
    "github.com/organization/service-seed/packages/stats"
)
```

## 🔧 Development Commands

```bash
# Build the application
make build

# Run with live reload
make dev

# Run tests
make test

# Build Docker image
make docker-build
```

This boilerplate provides a solid foundation for building scalable, maintainable Go applications. Always extend rather than replace, and follow the established patterns for consistency and reliability.
