package invoice

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	domainInvoice "github.com/your-org/jvairv2/pkg/domain/invoice"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Delete maneja la solicitud de eliminación de una factura (soft delete)
// @Summary Eliminar factura
// @Description Elimina una factura (soft delete)
// @Tags Invoices
// @Accept json
// @Produce json
// @Param id path int true "ID de la factura"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/invoices/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := h.useCase.Delete(r.Context(), id); err != nil {
		if err == domainInvoice.ErrInvoiceNotFound {
			response.Error(w, http.StatusNotFound, "Factura no encontrada")
			return
		}
		slog.ErrorContext(r.Context(), "Failed to delete invoice",
			slog.Int64("id", id),
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Error al eliminar factura")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
