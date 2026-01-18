package customer

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/your-org/jvairv2/pkg/rest/response"
)

// List godoc
// @Summary List customers
// @Description Get a paginated list of customers with optional filters
// @Tags customers
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Param search query string false "Search term"
// @Param workflow_id query int false "Filter by workflow ID"
// @Success 200 {object} response.PaginatedResponse{items=[]CustomerResponse}
// @Failure 500 {object} response.ErrorResponse
// @Router /customers [get]
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

	customers, total, err := h.useCase.List(r.Context(), filters, page, pageSize)
	if err != nil {
		// Errores de validación de negocio (workflow inválido o inactivo)
		if err.Error() == "invalid workflow_id" || err.Error() == "workflow is not active" {
			slog.WarnContext(r.Context(), "Invalid filter parameters",
				slog.String("error", err.Error()))
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		// Otros errores internos
		slog.ErrorContext(r.Context(), "Failed to list customers",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Failed to list customers")
		return
	}

	items := make([]CustomerResponse, len(customers))
	for i, c := range customers {
		items[i] = toCustomerResponse(c)
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
