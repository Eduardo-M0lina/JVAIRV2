package invoice

import (
	"log/slog"
	"net/http"
	"strconv"

	domainInvoice "github.com/angumol/jvairv2/pkg/domain/invoice"
	"github.com/angumol/jvairv2/pkg/rest/response"
	"github.com/go-chi/chi/v5"
)

// Get maneja la solicitud de obtención de una factura por ID
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
