package invoice_payment

import (
	"log/slog"
	"net/http"
	"strconv"

	domainPayment "github.com/angumol/jvairv2/pkg/domain/invoice_payment"
	"github.com/angumol/jvairv2/pkg/rest/response"
	"github.com/go-chi/chi/v5"
)

// Get maneja la solicitud de obtención de un pago por ID
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
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

	payment, err := h.useCase.GetByID(r.Context(), invoiceID, id)
	if err != nil {
		if err == domainPayment.ErrPaymentNotFound {
			response.Error(w, http.StatusNotFound, "Pago no encontrado")
			return
		}
		slog.ErrorContext(r.Context(), "Failed to get invoice payment",
			slog.Int64("id", id),
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Error al obtener pago")
		return
	}

	response.JSON(w, http.StatusOK, toPaymentResponse(payment))
}
