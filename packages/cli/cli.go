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
	var rootCmd = &cobra.Command{
		Use:   "service-seed",
		Short: "Service Seed - A production-ready Go application template",
	}
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

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

	rootCmd.AddCommand(cmdAgent)

	return rootCmd
}
