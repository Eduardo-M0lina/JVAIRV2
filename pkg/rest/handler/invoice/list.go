package invoice

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/angumol/jvairv2/pkg/rest/response"
)

// List maneja la solicitud de listado de facturas
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize < 1 {
		pageSize = 10
	}

	filters := parseFilters(r)

	invoices, total, err := h.useCase.List(r.Context(), filters, page, pageSize)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to list invoices",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Error al listar facturas")
		return
	}

	items := make([]InvoiceResponse, len(invoices))
	for i, inv := range invoices {
		items[i] = toInvoiceResponse(inv)
	}

	totalPages := (total + pageSize - 1) / pageSize

	response.JSON(w, http.StatusOK, response.PaginatedResponse{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	})
}
