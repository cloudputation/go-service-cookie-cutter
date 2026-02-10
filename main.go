package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cloudputation/service-seed/packages/cli"
	"github.com/cloudputation/service-seed/packages/config"
	log "github.com/cloudputation/service-seed/packages/logger"
	"github.com/cloudputation/service-seed/packages/stats"
)

func main() {
	fmt.Printf("INFO: Starting service-seed agent..\n\n")

	// Load main configuration file
	err := config.LoadConfiguration()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logging system first (before other components that may use it)
	logOpts := &log.LoggerOptions{}

	// Initialize OTLP log export if enabled
	if config.AppConfig.Telemetry != nil &&
		config.AppConfig.Telemetry.Logs != nil &&
		config.AppConfig.Telemetry.Logs.Enabled {

		// Note: OTLP log export will be implemented when log.InitOTLPLogs is added
		// For now, we configure the logger without OTLP support
		fmt.Printf("INFO: OTLP logs configured but not yet implemented (endpoint: %s)\n",
			config.AppConfig.Telemetry.Logs.Endpoint)
	}

	// Determine log level (defaulting to "info" if not configured)
	logLevel := "info"
	// Note: Add log_level field to config.Configuration when ready

	err = log.InitLoggerWithOptions(config.AppConfig.LogDir, logLevel, logOpts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logs: %v\n", err)
		os.Exit(1)
	}
	defer log.CloseLogger()

	// Graceful shutdown for OTLP logs (placeholder for future implementation)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// Note: Implement log.ShutdownOTLPLogs when OTLP adapter is added
		_ = ctx
	}()

	// Initialize server metrics
	err = stats.InitMetrics()
	if err != nil {
		log.Fatal("Failed to initialize metrics service: %v", err)
	}

	// Graceful shutdown for metrics (placeholder for OTLP exports)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		// Note: Implement stats.Shutdown when OTLP metric export is added
		_ = ctx
	}()

	// Initialize OTLP traces if enabled (placeholder for future implementation)
	if config.AppConfig.Telemetry != nil &&
		config.AppConfig.Telemetry.Traces != nil &&
		config.AppConfig.Telemetry.Traces.Enabled {

		// Note: Implement stats.InitTraces when OTLP trace export is added
		fmt.Printf("INFO: OTLP traces configured but not yet implemented (endpoint: %s, sampling: %.2f)\n",
			config.AppConfig.Telemetry.Traces.Endpoint,
			config.AppConfig.Telemetry.Traces.SamplingRate)
	}

	// Graceful shutdown for traces (placeholder for future implementation)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// Note: Implement stats.ShutdownTraces when OTLP trace export is added
		_ = ctx
	}()

	// Run CLI
	rootCmd := cli.SetupRootCommand()
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Error executing command: %v", err)
	}
}
