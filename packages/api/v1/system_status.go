package v1

import (
		"encoding/json"
		"fmt"
		"net/http"

		"github.com/organization/service-seed/packages/config"
		"github.com/organization/service-seed/packages/consul"
		"github.com/organization/service-seed/packages/stats"
		l "github.com/organization/service-seed/packages/logger"
)

type FileSystemStatus struct {
	Status string `json:"filesystem-status"`
}

type SystemStatusResponseBody struct {
	Message string `json:"message"`
}

func SystemStatusHandlerWrapper(w http.ResponseWriter, r *http.Request) {
	SystemStatusHandler(w, r)
}

func SystemStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
      err := http.StatusMethodNotAllowed
      l.Error("Received an invalid request method: %v", err)
      stats.ErrorCounter.Add(r.Context(), 1)
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
	}
	stats.SystemStatusEndpointCounter.Add(r.Context(), 1)

	filesystemDataPath := config.ConsulFileSystemDataDir
	statusPath := fmt.Sprintf("%s/status", filesystemDataPath)

	consulFileSystemStatus, err := consul.ConsulStoreGet(statusPath)
	if err != nil {
			l.Error("Failed to fetch filesystem state: "+err.Error(), http.StatusInternalServerError)
			stats.ErrorCounter.Add(r.Context(), 1)
			http.Error(w, "Failed to fetch filesystem state: "+err.Error(), http.StatusInternalServerError)
			return
	}

	consulFileSystemStatusBytes, err := json.Marshal(consulFileSystemStatus)
	if err != nil {
			l.Error("Error marshaling map to JSON: %v", err)
			stats.ErrorCounter.Add(r.Context(), 1)
			http.Error(w, "Failed to marshal map to JSON", http.StatusInternalServerError)
			return
	}

	var status FileSystemStatus
	err = json.Unmarshal(consulFileSystemStatusBytes, &status)
	if err != nil {
			l.Error("Error unmarshaling JSON: %v", err)
			stats.ErrorCounter.Add(r.Context(), 1)
			http.Error(w, "Failed to unmarshal JSON", http.StatusInternalServerError)
			return
	}


	w.Header().Set("Content-Type", "text/plain")
  fmt.Fprintf(w, status.Status)
}
