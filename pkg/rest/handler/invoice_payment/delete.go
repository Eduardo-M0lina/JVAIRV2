package invoice_payment

import (
	"log/slog"
	"net/http"
	"strconv"

	domainPayment "github.com/angumol/jvairv2/pkg/domain/invoice_payment"
	"github.com/angumol/jvairv2/pkg/rest/response"
	"github.com/go-chi/chi/v5"
)

// Delete maneja la solicitud de eliminación de un pago (soft delete)
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
		if err == domainPayment.ErrStripePaymentImmutable {
			response.Error(w, http.StatusForbidden, "Los pagos de Stripe no pueden ser eliminados")
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
