package v1

import (
    "net/http"
    "io/ioutil"

    l "github.com/organization/service-seed/packages/logger"
    "github.com/organization/service-seed/packages/stats"
)


func HealthHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
      err := http.StatusMethodNotAllowed
      l.ContextError(err, "HealthHandler: invalid request method")
      stats.ErrorCounter.Add(r.Context(), 1)
      http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
      return
  }
  stats.HealthEndpointCounter.Add(r.Context(), 1)


  content, err := ioutil.ReadFile("./VERSION")
  if err != nil {
      l.ContextError(err, "HealthHandler: failed to read API version")
      http.Error(w, "Failed to read API version", http.StatusInternalServerError)
      return
  }

  response := string(content) + "OK\n"
  w.Write([]byte(response))
}
