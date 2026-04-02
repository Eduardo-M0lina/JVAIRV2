package invoice

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	handler "github.com/your-org/jvairv2/pkg/rest/handler"
)

// SendEmailRequest representa la solicitud para enviar un email de invoice
type SendEmailRequest struct {
	Email string `json:"email"`
}

// SendEmail maneja el envío de email de invoice
// @Summary Send invoice email
// @Description Send invoice email to specified recipients
// @Tags invoices
// @Accept json
// @Produce json
// @Param id path int true "Invoice ID"
// @Param request body SendEmailRequest true "Email request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /invoices/{id}/email [post]
func (h *Handler) SendEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Obtener ID del invoice
	invoiceIDStr := chi.URLParam(r, "id")
	invoiceID, err := strconv.ParseInt(invoiceIDStr, 10, 64)
	if err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "ID de invoice inválido")
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
	err = h.emailService.SendInvoiceEmail(ctx, invoiceID, emails)
	if err != nil {
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Responder con éxito
	handler.RespondWithSuccess(w, "Invoice emailed to "+req.Email)
}
