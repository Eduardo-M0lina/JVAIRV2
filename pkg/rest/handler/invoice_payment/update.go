package invoice_payment

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	domainPayment "github.com/your-org/jvairv2/pkg/domain/invoice_payment"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Update maneja la solicitud de actualización de un pago
// @Summary Actualizar pago de factura
// @Description Actualiza un pago existente
// @Tags Invoice Payments
// @Accept json
// @Produce json
// @Param invoiceId path int true "ID de la factura"
// @Param id path int true "ID del pago"
// @Param payment body UpdatePaymentRequest true "Datos a actualizar"
// @Success 200 {object} PaymentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/invoices/{invoiceId}/payments/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	invoiceID, err := parseInvoiceID(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de factura inválido")
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de pago inválido")
		return
	}

	var req UpdatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.WarnContext(r.Context(), "Invalid request body",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	// Obtener pago existente para hacer merge
	existing, err := h.useCase.GetByID(r.Context(), invoiceID, id)
	if err != nil {
		if err == domainPayment.ErrPaymentNotFound {
			response.Error(w, http.StatusNotFound, "Pago no encontrado")
			return
		}
		slog.ErrorContext(r.Context(), "Failed to get payment for update",
			slog.Int64("id", id),
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Error al obtener pago")
		return
	}

	// Merge campos
	payment := &domainPayment.InvoicePayment{
		ID:               id,
		InvoiceID:        invoiceID,
		PaymentProcessor: existing.PaymentProcessor,
		PaymentID:        existing.PaymentID,
		Amount:           existing.Amount,
		Notes:            existing.Notes,
	}

	if req.PaymentProcessor != nil {
		payment.PaymentProcessor = *req.PaymentProcessor
	}
	if req.PaymentID != nil {
		payment.PaymentID = *req.PaymentID
	}
	if req.Amount != nil {
		payment.Amount = *req.Amount
	}
	if req.Notes != nil {
		payment.Notes = *req.Notes
	}

	if err := h.useCase.Update(r.Context(), payment); err != nil {
		switch err {
		case domainPayment.ErrPaymentNotFound:
			response.Error(w, http.StatusNotFound, "Pago no encontrado")
		default:
			if err.Error() == "id is required" {
				response.Error(w, http.StatusBadRequest, err.Error())
			} else {
				response.Error(w, http.StatusInternalServerError, "Error al actualizar pago")
			}
		}
		return
	}

	// Re-obtener el pago actualizado
	updated, err := h.useCase.GetByID(r.Context(), invoiceID, id)
	if err != nil {
		response.JSON(w, http.StatusOK, toPaymentResponse(payment))
		return
	}

	response.JSON(w, http.StatusOK, toPaymentResponse(updated))
}
