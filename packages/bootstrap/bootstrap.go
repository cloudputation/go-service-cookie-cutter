package bootstrap

import (
    "fmt"
    "os"

    "github.com/organization/service-seed/packages/config"
    "github.com/organization/service-seed/packages/consul"
    l "github.com/organization/service-seed/packages/logger"
    "github.com/organization/service-seed/packages/stats"
)

func BootstrapFileSystem() error {
  l.Info("Starting Service FileSystem agent.. Bootstrapping filesystem.")
  dataDir := config.AppConfig.DataDir
  rootDir := config.RootDir
  l.Info("Loaded configuration file: %s", config.ConfigPath)

  err = consul.InitConsul()
  if err != nil {
      return fmt.Errorf("Could not initialize Consul: %v", err)
  }

  err = consul.BootstrapConsul()
  if err != nil {
      return fmt.Errorf("Could not bootstrap filesystem on Consul: %v", err)
  }

  l.Info("Refreshing filesystem state.")
  err = stats.GenerateState()
  if err != nil {
      return fmt.Errorf("Failed to generate filesystem state: %v", err)
  }
  l.Info("FileSystem state created successfully!")
  l.Info("FileSystem bootstrapping done!")


  return nil
}
