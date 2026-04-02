package api

import (
	"encoding/json"
	"net/http"
	"time"
)

// HealthStatus holds the service health information.
type HealthStatus struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Uptime    string `json:"uptime"`
	Version   string `json:"version"`
}

// HealthHandler provides HTTP health check endpoints.
type HealthHandler struct {
	startTime time.Time
	version   string
}

// NewHealthHandler creates a new HealthHandler with the given version.
func NewHealthHandler(version string) *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
		version:   version,
	}
}

// ServeHTTP handles GET requests and returns service health as JSON.
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	uptime := time.Since(h.startTime).Round(time.Second)
	status := HealthStatus{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Uptime:    uptime.String(),
		Version:   h.version,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}

// IsHealthy returns true if the service is in a healthy state.
func (h *HealthHandler) IsHealthy() bool {
	return true
}
