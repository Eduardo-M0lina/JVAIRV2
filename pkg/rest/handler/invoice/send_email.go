package invoice

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
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
		http.Error(w, `{"error":"ID de invoice inválido"}`, http.StatusBadRequest)
		return
	}

	// Parsear request
	var req SendEmailRequest
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
	err = h.emailService.SendInvoiceEmail(ctx, invoiceID, emails)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Responder con éxito
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Invoice emailed to " + req.Email,
	})
}
