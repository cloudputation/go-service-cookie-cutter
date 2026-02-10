package v1

import (
    "net/http"

    log "github.com/cloudputation/service-seed/packages/logger"
    "github.com/cloudputation/service-seed/packages/stats"
)


func HealthHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
      log.Error("HealthHandler: invalid request method")
      stats.ErrorCounter.Add(r.Context(), 1)
      http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
      return
  }
  stats.HealthEndpointCounter.Add(r.Context(), 1)

  w.Header().Set("Content-Type", "text/plain")
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("OK\n"))
}
