package invoice

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	domainInvoice "github.com/your-org/jvairv2/pkg/domain/invoice"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Get maneja la solicitud de obtención de una factura por ID
// @Summary Obtener factura
// @Description Obtiene una factura por su ID, incluyendo balance calculado
// @Tags Invoices
// @Accept json
// @Produce json
// @Param id path int true "ID de la factura"
// @Success 200 {object} InvoiceResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/invoices/{id} [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	inv, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if err == domainInvoice.ErrInvoiceNotFound {
			response.Error(w, http.StatusNotFound, "Factura no encontrada")
			return
		}
		slog.ErrorContext(r.Context(), "Failed to get invoice",
			slog.Int64("id", id),
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Error al obtener factura")
		return
	}

	response.JSON(w, http.StatusOK, toInvoiceResponse(inv))
}
