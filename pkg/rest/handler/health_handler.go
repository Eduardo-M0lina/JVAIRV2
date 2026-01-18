package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

// DBChecker define la interfaz para verificar el estado de la base de datos
type DBChecker interface {
	Ping() error
}

// HealthHandler maneja las solicitudes de health check
type HealthHandler struct {
	DBChecker DBChecker
}

// HealthResponse representa la respuesta del health check
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	DBStatus  string    `json:"dbStatus"`
}

// NewHealthHandler crea un nuevo handler para health check
func NewHealthHandler(dbChecker DBChecker) *HealthHandler {
	return &HealthHandler{
		DBChecker: dbChecker,
	}
}

// Check maneja la solicitud GET para verificar el estado del servicio
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		DBStatus:  "ok",
	}

	// Verificar conexión a la base de datos
	if h.DBChecker != nil {
		if err := h.DBChecker.Ping(); err != nil {
			response.Status = "error"
			response.DBStatus = "error: " + err.Error()
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Error al codificar respuesta de health check",
			"error", err,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
