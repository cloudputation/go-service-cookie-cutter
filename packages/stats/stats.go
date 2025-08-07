package stats

import (
    "fmt"
    "io/ioutil"
    "encoding/json"

    "go.opentelemetry.io/otel/exporters/prometheus"
    api "go.opentelemetry.io/otel/metric"
    "go.opentelemetry.io/otel/sdk/metric"

    "github.com/organization/service-seed/packages/config"
    "github.com/organization/service-seed/packages/consul"
)


type FileSystemInfo struct {
  FileSystemState   string   `json:"filesystem-state"`
  Services          []string `json:"services"`
}

const meterName = "CFS.Metrics"
var (
    ErrorCounter                 api.Int64Counter
    HealthEndpointCounter        api.Int64Counter
    ServiceStatusEndpointCounter api.Int64Counter
    SystemMetricsEndpointCounter api.Int64Counter
)


func InitMetrics() error {
    exporter, err := prometheus.New()
    if err != nil {
        return fmt.Errorf("Failed to initialize prometheus client: %v", err)
    }
    provider := metric.NewMeterProvider(metric.WithReader(exporter))
    meter := provider.Meter(meterName)

    ErrorCounter, err = meter.Int64Counter(
        "agent_errors",
        api.WithDescription("Counts the number of errors during agent runtime"),
    )
    if err != nil {
        return fmt.Errorf("Failed to initialize the error counter: %v", err)
    }

    HealthEndpointCounter, err = meter.Int64Counter(
        "health_endpoint_hits",
        api.WithDescription("Counts the number of hits to the /health endpoint"),
    )
    if err != nil {
        return fmt.Errorf("Failed to initialize /health endpoint counter: %v", err)
    }

    ServiceStatusEndpointCounter, err = meter.Int64Counter(
        "service_status_endpoint_hits",
        api.WithDescription("Counts the number of hits to the /service/status endpoint"),
    )
    if err != nil {
        return fmt.Errorf("Failed to initialize /service/status endpoint counter: %v", err)
    }

    SystemMetricsEndpointCounter, err = meter.Int64Counter(
        "system_metrics_endpoint_hits",
        api.WithDescription("Counts the number of hits to the /system/metrics endpoint"),
    )
    if err != nil {
        return fmt.Errorf("Failed to initialize /system/metrics endpoint counter: %v", err)
    }


    return nil
}

func GenerateState() error {
  dirName := config.RootDir + "/" + config.AppConfig.DataDir + "/services"
  jsonData, err := GetFileSystemInfo(dirName)
	if err != nil {
      return fmt.Errorf("Failed to get filesystem info: %v", err)
	}

	keyPath := config.ConsulServiceSummaryDataDir
	err = consul.ConsulStorePut(jsonData, keyPath)
	if err != nil {
      return fmt.Errorf("Failed to put data in Consul KV: %v", err)
	}


  return nil
}

func GetFileSystemInfo(dirPath string) (string, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
      return "", fmt.Errorf("Failed to read filesystem files in %s: %v", dirPath, err)
	}

	var dirNames []string
	for _, file := range files {
      if file.IsDir() {
          dirNames = append(dirNames, file.Name())
      }
	}

	filesystemInfo := FileSystemInfo{
		FileSystemState: "running",
		Services:     dirNames,
	}

	jsonData, err := json.MarshalIndent(filesystemInfo, "", "  ")
	if err != nil {
      return "", fmt.Errorf("Failed to process filesystem info: %v", err)
	}


	return string(jsonData), nil
}
