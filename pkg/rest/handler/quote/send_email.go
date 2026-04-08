package quote

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	handler "github.com/your-org/jvairv2/pkg/rest/handler"
)

// SendEmailRequest representa la solicitud para enviar un email de quote
type SendEmailRequest struct {
	Email string `json:"email"`
}

// SendEmail maneja el envío de email de quote
func (h *Handler) SendEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Obtener ID del quote
	quoteIDStr := chi.URLParam(r, "id")
	quoteID, err := strconv.ParseInt(quoteIDStr, 10, 64)
	if err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "ID de quote inválido")
		return
	}

	// Parsear request
	var req SendEmailRequest
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
	err = h.emailService.SendQuoteEmail(ctx, quoteID, emails)
	if err != nil {
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Responder con éxito
	handler.RespondWithSuccess(w, "Quote emailed to "+req.Email)
}
