package invoice

import (
	"encoding/json"
	"log/slog"
	"net/http"

	domainInvoice "github.com/your-org/jvairv2/pkg/domain/invoice"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Create maneja la solicitud de creación de una factura
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.WarnContext(r.Context(), "Invalid request body",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	inv := &domainInvoice.Invoice{
		JobID:         req.JobID,
		InvoiceNumber: req.InvoiceNumber,
		Total:         req.Total,
		Description:   req.Description,
		Notes:         req.Notes,
	}

	if req.AllowOnlinePayments != nil {
		inv.AllowOnlinePayments = *req.AllowOnlinePayments
	}

	if err := h.useCase.Create(r.Context(), inv); err != nil {
		switch err {
		case domainInvoice.ErrInvalidJob:
			response.Error(w, http.StatusBadRequest, err.Error())
		default:
			if err.Error() == "job_id is required" ||
				err.Error() == "invoice_number is required" {
				response.Error(w, http.StatusBadRequest, err.Error())
			} else {
				response.Error(w, http.StatusInternalServerError, "Error al crear factura")
			}
		}
		return
	}

	response.JSON(w, http.StatusCreated, toInvoiceResponse(inv))
}
