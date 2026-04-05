package invoice_payment

import (
	"encoding/json"
	"log/slog"
	"net/http"

	domainPayment "github.com/your-org/jvairv2/pkg/domain/invoice_payment"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Create maneja la solicitud de creación de un pago de factura
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	invoiceID, err := parseInvoiceID(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de factura inválido")
		return
	}

	var req CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.WarnContext(r.Context(), "Invalid request body",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	payment := &domainPayment.InvoicePayment{
		InvoiceID:        invoiceID,
		PaymentProcessor: req.PaymentProcessor,
		PaymentID:        req.PaymentID,
		Amount:           req.Amount,
		Notes:            req.Notes,
	}

	if err := h.useCase.Create(r.Context(), payment); err != nil {
		switch err {
		case domainPayment.ErrInvalidInvoice:
			response.Error(w, http.StatusNotFound, "Factura no encontrada")
		default:
			if err.Error() == "invoice_id is required" ||
				err.Error() == "payment_processor is required" ||
				err.Error() == "payment_id is required" ||
				err.Error() == "amount must be greater than 0" {
				response.Error(w, http.StatusBadRequest, err.Error())
			} else {
				response.Error(w, http.StatusInternalServerError, "Error al crear pago")
			}
		}
		return
	}

	response.JSON(w, http.StatusCreated, toPaymentResponse(payment))
}
