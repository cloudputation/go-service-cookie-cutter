package api

import (
    "fmt"
    "net/http"

    "github.com/prometheus/client_golang/prometheus/promhttp"

    "github.com/organization/service-seed/packages/config"
    l "github.com/organization/service-seed/packages/logger"
    "github.com/organization/service-seed/packages/api/v1"
)


const MaxWorkers = 10

func StartServer() {
  serverPort := fmt.Sprintf(":%s", config.AppConfig.Server.ServerPort)
  l.Info("Starting server on port %s", serverPort)

  http.HandleFunc("/v1/health", v1.HealthHandler)
  http.HandleFunc("/v1/system/status", v1.SystemStatusHandlerWrapper)
  http.Handle("/v1/system/metrics", promhttp.Handler())

  err := http.ListenAndServe(serverPort, nil)
  if err != nil {
      l.Fatal("HTTP server failed to start: %v", err)
  }
}
