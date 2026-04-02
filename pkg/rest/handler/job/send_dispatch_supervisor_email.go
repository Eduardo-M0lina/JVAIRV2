package job

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	handler "github.com/your-org/jvairv2/pkg/rest/handler"
)

// SendDispatchSupervisorEmailRequest representa la solicitud para enviar email a supervisores
type SendDispatchSupervisorEmailRequest struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// SendDispatchSupervisorEmail maneja el envío de email a supervisores
// @Summary Send dispatch supervisor email
// @Description Send custom email to supervisors
// @Tags jobs
// @Accept json
// @Produce json
// @Param id path int true "Job ID"
// @Param request body SendDispatchSupervisorEmailRequest true "Email request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /jobs/{id}/dispatch-supervisor-email [post]
func (h *Handler) SendDispatchSupervisorEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Obtener ID del job
	jobIDStr := chi.URLParam(r, "id")
	jobID, err := strconv.ParseInt(jobIDStr, 10, 64)
	if err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "ID de job inválido")
		return
	}

	// Parsear request
	var req SendDispatchSupervisorEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Request inválido")
		return
	}

	// Validar campos requeridos
	if req.Email == "" {
		handler.RespondWithError(w, http.StatusBadRequest, "Email es requerido")
		return
	}
	if req.Subject == "" {
		handler.RespondWithError(w, http.StatusBadRequest, "Subject es requerido")
		return
	}
	if req.Body == "" {
		handler.RespondWithError(w, http.StatusBadRequest, "Body es requerido")
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
	err = h.emailService.SendDispatchSupervisorEmail(ctx, jobID, req.Subject, req.Body, emails)
	if err != nil {
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Responder con éxito
	handler.RespondWithSuccess(w, "Dispatch supervisor email sent to "+req.Email)
}
