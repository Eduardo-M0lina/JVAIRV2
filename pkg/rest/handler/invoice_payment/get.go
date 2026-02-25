package invoice_payment

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	domainPayment "github.com/your-org/jvairv2/pkg/domain/invoice_payment"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Get maneja la solicitud de obtención de un pago por ID
// @Summary Obtener pago de factura
// @Description Obtiene un pago por su ID
// @Tags Invoice Payments
// @Accept json
// @Produce json
// @Param invoiceId path int true "ID de la factura"
// @Param id path int true "ID del pago"
// @Success 200 {object} PaymentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/invoices/{invoiceId}/payments/{id} [get]
// @Security BearerAuth
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
