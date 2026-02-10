package bootstrap

import (
    "fmt"
    "os"

    "github.com/organization/service-seed/packages/config"
    l "github.com/organization/service-seed/packages/logger"
)

func BootstrapFileSystem() error {
  l.Info("Starting Service agent.. Bootstrapping filesystem.")
  dataDir := config.AppConfig.DataDir
  rootDir := config.RootDir
  l.Info("Loaded configuration file: %s", config.ConfigPath)

  // Ensure data directory exists
  dataDirPath := rootDir + "/" + dataDir
  err := os.MkdirAll(dataDirPath, 0755)
  if err != nil {
      return fmt.Errorf("Failed to create data directory: %v", err)
  }

  l.Info("Data directory initialized at: %s", dataDirPath)
  l.Info("FileSystem bootstrapping done!")

  return nil
}
