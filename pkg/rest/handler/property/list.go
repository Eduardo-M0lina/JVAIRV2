package property

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/your-org/jvairv2/pkg/rest/response"
)

// List godoc
// @Summary List properties
// @Description Get a paginated list of properties with optional filters
// @Tags Properties
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Param search query string false "Search term"
// @Param customerId query int false "Filter by customer ID"
// @Success 200 {object} response.PaginatedResponse{items=[]PropertyResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/properties [get]
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

	properties, total, err := h.useCase.List(r.Context(), filters, page, pageSize)
	if err != nil {
		if err.Error() == "invalid customer_id" || err.Error() == "customer is deleted" {
			slog.WarnContext(r.Context(), "Invalid filter parameters",
				slog.String("error", err.Error()))
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		slog.ErrorContext(r.Context(), "Failed to list properties",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Failed to list properties")
		return
	}

	items := make([]PropertyResponse, len(properties))
	for i, p := range properties {
		items[i] = toPropertyResponse(p)
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
