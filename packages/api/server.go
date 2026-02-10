package api

import (
    "fmt"
    "net/http"

    "github.com/prometheus/client_golang/prometheus/promhttp"

    "github.com/cloudputation/service-seed/packages/config"
    log "github.com/cloudputation/service-seed/packages/logger"
    "github.com/cloudputation/service-seed/packages/api/v1"
)


const MaxWorkers = 10

func StartServer() {
  serverPort := fmt.Sprintf(":%s", config.AppConfig.Server.ServerPort)
  log.Info("Starting server on port %s", serverPort)

  http.HandleFunc("/v1/health", v1.HealthHandler)
  http.Handle("/v1/system/metrics", promhttp.Handler())

  err := http.ListenAndServe(serverPort, nil)
  if err != nil {
      log.Fatal("HTTP server failed to start: %v", err)
  }
}
