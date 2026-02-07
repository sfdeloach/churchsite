package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sfdeloach/churchsite/internal/database"
)

// HealthHandler handles health check endpoints.
type HealthHandler struct {
	db *database.DB
}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler(db *database.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Liveness responds with a simple OK status for container health checks.
func (h *HealthHandler) Liveness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Readiness checks that both PostgreSQL and Redis are reachable.
func (h *HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pgStatus := "ok"
	if err := h.db.PingPostgres(); err != nil {
		pgStatus = "error"
	}

	redisStatus := "ok"
	if err := h.db.PingRedis(); err != nil {
		redisStatus = "error"
	}

	status := "ok"
	code := http.StatusOK
	if pgStatus != "ok" || redisStatus != "ok" {
		status = "degraded"
		code = http.StatusServiceUnavailable
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"status":   status,
		"postgres": pgStatus,
		"redis":    redisStatus,
	})
}
