package invoice_payment

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	domainPayment "github.com/your-org/jvairv2/pkg/domain/invoice_payment"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Delete maneja la solicitud de eliminación de un pago (soft delete)
// @Summary Eliminar pago de factura
// @Description Elimina un pago (soft delete)
// @Tags Invoice Payments
// @Accept json
// @Produce json
// @Param invoiceId path int true "ID de la factura"
// @Param id path int true "ID del pago"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/invoices/{invoiceId}/payments/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
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

	if err := h.useCase.Delete(r.Context(), invoiceID, id); err != nil {
		if err == domainPayment.ErrPaymentNotFound {
			response.Error(w, http.StatusNotFound, "Pago no encontrado")
			return
		}
		slog.ErrorContext(r.Context(), "Failed to delete invoice payment",
			slog.Int64("id", id),
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Error al eliminar pago")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
