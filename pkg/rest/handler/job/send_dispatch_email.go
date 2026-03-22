package job

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	handler "github.com/your-org/jvairv2/pkg/rest/handler"
)

// SendDispatchEmailRequest representa la solicitud para enviar email de dispatch
type SendDispatchEmailRequest struct {
	Email string `json:"email"`
}

// SendDispatchEmail envía un email de dispatch del job a técnicos
// @Summary Enviar email de dispatch
// @Description Envía un email de dispatch del job a uno o más técnicos
// @Tags jobs
// @Accept json
// @Produce json
// @Param id path int true "Job ID"
// @Param request body SendDispatchEmailRequest true "Email request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /jobs/{id}/dispatch-email [post]
func (h *Handler) SendDispatchEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Obtener ID del job
	jobIDStr := chi.URLParam(r, "id")
	jobID, err := strconv.ParseInt(jobIDStr, 10, 64)
	if err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "ID de job inválido")
		return
	}

	// Parsear request
	var req SendDispatchEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Request inválido")
		return
	}

	// Validar emails
	if req.Email == "" {
		handler.RespondWithError(w, http.StatusBadRequest, "Email es requerido")
		return
	}

	// Separar emails por coma y limpiar espacios
	emails := strings.Split(req.Email, ",")
	for i, email := range emails {
		emails[i] = strings.TrimSpace(email)
	}

	// Verificar que el servicio de email esté disponible
	if h.emailService == nil {
		handler.RespondWithError(w, http.StatusInternalServerError, "Servicio de email no configurado")
		return
	}

	// Enviar email usando el servicio de dominio
	err = h.emailService.SendDispatchEmail(ctx, jobID, emails)
	if err != nil {
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Registrar actividad
	h.logActivity(ctx, jobID, "email_sent", fmt.Sprintf("Dispatch email sent to: %s", req.Email))

	// Responder
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Job emailed to %s", req.Email),
	})
}
