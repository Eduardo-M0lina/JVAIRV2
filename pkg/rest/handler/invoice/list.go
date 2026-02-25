package invoice

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/your-org/jvairv2/pkg/rest/response"
)

// List maneja la solicitud de listado de facturas
// @Summary Listar facturas
// @Description Obtiene una lista paginada de facturas con filtros opcionales
// @Tags Invoices
// @Accept json
// @Produce json
// @Param page query int false "Número de página" default(1)
// @Param pageSize query int false "Tamaño de página" default(10)
// @Param search query string false "Búsqueda en invoice_number, work_order, property, customer"
// @Param jobId query int false "Filtrar por job"
// @Param status query string false "Filtrar por status: paid, unpaid"
// @Param sort query string false "Campo de ordenamiento: invoice_number, total, balance, created_at"
// @Param direction query string false "Dirección de ordenamiento: asc, desc" default(desc)
// @Success 200 {object} response.PaginatedResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/invoices [get]
// @Security BearerAuth
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
