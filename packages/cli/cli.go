package cli

import (
	"context"
	"github.com/spf13/cobra"

	"github.com/organization/service-seed/packages/bootstrap"
	l "github.com/organization/service-seed/packages/logger"
	"github.com/organization/service-seed/packages/api"
	"github.com/organization/service-seed/packages/stats"
)


func SetupRootCommand() *cobra.Command {
	var serviceName  string
	var serviceFiles []string
	var serviceNames []string

	var rootCmd = &cobra.Command{
		Use:   "service-seed",
		Short: "Service Seed.",
	}
	rootCmd.CompletionOptions.HiddenDefaultCmd = true


	var cmdConfig = &cobra.Command{
		Use:   "config",
		Short: "Check the configuration",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			err := CheckConfig()
			if err != nil {
				l.Error("Failed to parse configuration file: %v", err)
				stats.ErrorCounter.Add(context.Background(), 1)
			}
		},
	}

	var cmdAgent = &cobra.Command{
		Use:   "agent",
		Short: "Start the service agent",
		Run: func(cmd *cobra.Command, args []string) {
			err := bootstrap.BootstrapFileSystem()
			if err != nil {
				l.Fatal("Failed to bootstrap the filesystem: %v", err)
			}
			api.StartServer()
		},
	}

	var cmdSystem = &cobra.Command{
		Use:   "system",
		Short: "Commands related to system operations",
	}

	var cmdSystemStatus = &cobra.Command{
		Use:   "status",
		Short: "Check the system status",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			err := GetFileSystemStatus()
			if err != nil {
				l.Error("Failed to check server status %v", err)
				stats.ErrorCounter.Add(context.Background(), 1)
			}
		},
	}
	cmdSystem.AddCommand(cmdSystemStatus)


	var commands = []*cobra.Command{
		cmdConfig,
		cmdAgent,
		cmdSystem,
	}

	for _, cmd := range commands {
		rootCmd.AddCommand(cmd)
	}

  return rootCmd

}
