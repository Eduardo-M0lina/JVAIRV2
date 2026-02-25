package invoice

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	domainInvoice "github.com/your-org/jvairv2/pkg/domain/invoice"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Update maneja la solicitud de actualización de una factura
// @Summary Actualizar factura
// @Description Actualiza una factura existente
// @Tags Invoices
// @Accept json
// @Produce json
// @Param id path int true "ID de la factura"
// @Param invoice body UpdateInvoiceRequest true "Datos a actualizar"
// @Success 200 {object} InvoiceResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/invoices/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var req UpdateInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.WarnContext(r.Context(), "Invalid request body",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	// Obtener factura existente para hacer merge
	existing, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if err == domainInvoice.ErrInvoiceNotFound {
			response.Error(w, http.StatusNotFound, "Factura no encontrada")
			return
		}
		slog.ErrorContext(r.Context(), "Failed to get invoice for update",
			slog.Int64("id", id),
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Error al obtener factura")
		return
	}

	// Merge campos
	inv := &domainInvoice.Invoice{
		ID:                  id,
		JobID:               existing.JobID,
		InvoiceNumber:       existing.InvoiceNumber,
		Total:               existing.Total,
		Description:         existing.Description,
		AllowOnlinePayments: existing.AllowOnlinePayments,
		Notes:               existing.Notes,
	}

	if req.JobID != nil {
		inv.JobID = *req.JobID
	}
	if req.InvoiceNumber != nil {
		inv.InvoiceNumber = *req.InvoiceNumber
	}
	if req.Total != nil {
		inv.Total = *req.Total
	}
	if req.Description != nil {
		inv.Description = req.Description
	}
	if req.AllowOnlinePayments != nil {
		inv.AllowOnlinePayments = *req.AllowOnlinePayments
	}
	if req.Notes != nil {
		inv.Notes = req.Notes
	}

	if err := h.useCase.Update(r.Context(), inv); err != nil {
		switch err {
		case domainInvoice.ErrInvoiceNotFound:
			response.Error(w, http.StatusNotFound, "Factura no encontrada")
		case domainInvoice.ErrInvalidJob:
			response.Error(w, http.StatusBadRequest, err.Error())
		default:
			if err.Error() == "id is required" {
				response.Error(w, http.StatusBadRequest, err.Error())
			} else {
				response.Error(w, http.StatusInternalServerError, "Error al actualizar factura")
			}
		}
		return
	}

	// Re-obtener para incluir balance actualizado
	updated, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		response.JSON(w, http.StatusOK, toInvoiceResponse(inv))
		return
	}

	response.JSON(w, http.StatusOK, toInvoiceResponse(updated))
}
