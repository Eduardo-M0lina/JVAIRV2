package user

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	handler "github.com/angumol/jvairv2/pkg/rest/handler"
	"github.com/go-chi/chi/v5"
)

// SendPayStubEmailRequest representa la solicitud para enviar email de paystub
type SendPayStubEmailRequest struct {
	Email string `json:"email"`
}

// SendPayStubEmail maneja el envío de email de recibo de pago
// @Summary Send paystub email
// @Description Send paystub/pay report email to specified recipients
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body SendPayStubEmailRequest true "Email request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/paystub-email [post]
func (h *Handler) SendPayStubEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Obtener ID del usuario
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "ID de usuario inválido")
		return
	}

	// Parsear request
	var req SendPayStubEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Request inválido")
		return
	}

	// Validar email
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
	err = h.emailService.SendPayStubEmail(ctx, userID, emails)
	if err != nil {
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Responder con éxito
	handler.RespondWithSuccess(w, "PayStub email sent to "+req.Email)
}
