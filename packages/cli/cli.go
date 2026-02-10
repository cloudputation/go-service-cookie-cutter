package cli

import (
	"github.com/spf13/cobra"

	"github.com/cloudputation/service-seed/packages/api"
	"github.com/cloudputation/service-seed/packages/bootstrap"
	log "github.com/cloudputation/service-seed/packages/logger"
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
				log.Fatal("Failed to bootstrap the filesystem: %v", err)
			}
			api.StartServer()
		},
	}

	rootCmd.AddCommand(cmdAgent)

	return rootCmd
}
