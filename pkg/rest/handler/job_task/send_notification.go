package job_task

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

// SendNotificationRequest representa la solicitud para enviar notificación de tarea
type SendNotificationRequest struct {
	Email string `json:"email"`
}

// SendNotification maneja el envío de notificación de tarea
// @Summary Send task notification email
// @Description Send task notification email to specified recipients
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param request body SendNotificationRequest true "Email request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id}/notification [post]
func (h *Handler) SendNotification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Obtener ID de la tarea
	taskIDStr := chi.URLParam(r, "id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"ID de tarea inválido"}`, http.StatusBadRequest)
		return
	}

	// Parsear request
	var req SendNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Request inválido"}`, http.StatusBadRequest)
		return
	}

	// Validar emails
	if req.Email == "" {
		http.Error(w, `{"error":"Email es requerido"}`, http.StatusBadRequest)
		return
	}

	// Separar emails por coma y limpiar espacios
	emails := strings.Split(req.Email, ",")
	for i, email := range emails {
		emails[i] = strings.TrimSpace(email)
	}

	// Verificar que el servicio de email esté disponible
	if h.emailService == nil {
		http.Error(w, `{"error":"Servicio de email no configurado"}`, http.StatusInternalServerError)
		return
	}

	// Enviar email usando el servicio de dominio
	err = h.emailService.SendTaskNotificationEmail(ctx, taskID, emails)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Responder con éxito
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Task notification emailed to " + req.Email,
	})
}
