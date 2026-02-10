package config

import (
    "os"
    "fmt"

    "github.com/spf13/viper"
    "github.com/hashicorp/hcl/v2/gohcl"
    "github.com/hashicorp/hcl/v2/hclparse"
)


type Configuration struct {
    LogDir      string      `hcl:"log_dir"`
    DataDir     string      `hcl:"data_dir"`
    Server      Server      `hcl:"server,block"`
    Telemetry   *Telemetry  `hcl:"telemetry,block"`
}

type Server struct {
    ServerPort    string `hcl:"port"`
    ServerAddress string `hcl:"address"`
}


var AppConfig Configuration
var ConfigPath string
var RootDir string
const MaxWorkers = 10


func LoadConfiguration() error {
  viper.SetDefault("ConfigPath", "/etc/service-seed/config.hcl")
  viper.BindEnv("ConfigPath", "SS_CONFIG_FILE_PATH")

  ConfigPath = viper.GetString("ConfigPath")

  var err error
  RootDir, err = os.Getwd()
  if err != nil {
      return fmt.Errorf("Failed to get service root directory: %v", err)
  }

  // Read the HCL file
  data, err := os.ReadFile(ConfigPath)
  if err != nil {
      return fmt.Errorf("Failed to read configuration file: %v", err)
  }

  // Parse the HCL file
  parser := hclparse.NewParser()
  file, diags := parser.ParseHCL(data, ConfigPath)
  if diags.HasErrors() {
      return fmt.Errorf("Failed to parse configuration: %v", diags)
  }

  // Populate the Config struct
  diags = gohcl.DecodeBody(file.Body, nil, &AppConfig)
  if diags.HasErrors() {
      return fmt.Errorf("Failed to apply configuration: %v", diags)
  }

  // Apply defaults for any missing optional values
  applyDefaults()

  return nil
}

func GetConfigPath() string {
  return ConfigPath
}

// DEFAULT CONFIGURATION SETTINGS
//
// applyDefaults sets default values for optional configuration fields
func applyDefaults() {
  applyTelemetryDefaults()
}
