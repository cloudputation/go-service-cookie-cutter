package v1

import (
		"encoding/json"
		"fmt"
		"net/http"

		"github.com/cloudputation/service-seed/packages/config"
		"github.com/cloudputation/service-seed/packages/stats"
		log "github.com/cloudputation/service-seed/packages/logger"
)

type SystemStatusResponse struct {
	Status  string `json:"status"`
	DataDir string `json:"data_dir"`
	LogDir  string `json:"log_dir"`
}

func SystemStatusHandlerWrapper(w http.ResponseWriter, r *http.Request) {
	SystemStatusHandler(w, r)
}

func SystemStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		err := http.StatusMethodNotAllowed
		log.Error("Received an invalid request method: %v", err)
		stats.ErrorCounter.Add(r.Context(), 1)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	stats.SystemStatusEndpointCounter.Add(r.Context(), 1)

	// Build simple system status response
	response := SystemStatusResponse{
		Status:  "running",
		DataDir: config.AppConfig.DataDir,
		LogDir:  config.AppConfig.LogDir,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Failed to encode system status response: %v", err)
		stats.ErrorCounter.Add(r.Context(), 1)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	log.Info("System status request completed successfully")
}
