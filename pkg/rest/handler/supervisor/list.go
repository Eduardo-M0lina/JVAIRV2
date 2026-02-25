package supervisor

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// List godoc
// @Summary List supervisors
// @Description Get a paginated list of supervisors with optional filters
// @Tags supervisors
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Param search query string false "Search term (name, phone, email)"
// @Param customerId query int false "Filter by customer ID"
// @Success 200 {object} response.PaginatedResponse{items=[]SupervisorResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /supervisors [get]
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

	supervisors, total, err := h.useCase.List(r.Context(), filters, page, pageSize)
	if err != nil {
		if err.Error() == "invalid customer_id" {
			slog.WarnContext(r.Context(), "Invalid filter parameters",
				slog.String("error", err.Error()))
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		slog.ErrorContext(r.Context(), "Failed to list supervisors",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Failed to list supervisors")
		return
	}

	items := make([]SupervisorResponse, len(supervisors))
	for i, s := range supervisors {
		items[i] = toSupervisorResponse(s)
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

// ListByCustomer godoc
// @Summary List supervisors by customer
// @Description Get a paginated list of supervisors for a specific customer
// @Tags supervisors,customers
// @Produce json
// @Param customerId path int true "Customer ID"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Param search query string false "Search term (name, phone, email)"
// @Success 200 {object} response.PaginatedResponse{items=[]SupervisorResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /customers/{customerId}/supervisors [get]
// @Security BearerAuth
func (h *Handler) ListByCustomer(w http.ResponseWriter, r *http.Request) {
	customerIDStr := chi.URLParam(r, "customerId")
	customerID, err := strconv.ParseInt(customerIDStr, 10, 64)
	if err != nil {
		slog.WarnContext(r.Context(), "Invalid customer ID",
			slog.String("id", customerIDStr))
		response.Error(w, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize < 1 {
		pageSize = 10
	}

	filters := map[string]interface{}{
		"customer_id": customerID,
	}

	if search := r.URL.Query().Get("search"); search != "" {
		filters["search"] = search
	}

	supervisors, total, err := h.useCase.List(r.Context(), filters, page, pageSize)
	if err != nil {
		if err.Error() == "invalid customer_id" {
			response.Error(w, http.StatusNotFound, "Customer not found")
			return
		}

		slog.ErrorContext(r.Context(), "Failed to list customer supervisors",
			slog.String("error", err.Error()),
			slog.Int64("customerId", customerID))
		response.Error(w, http.StatusInternalServerError, "Failed to list supervisors")
		return
	}

	items := make([]SupervisorResponse, len(supervisors))
	for i, s := range supervisors {
		items[i] = toSupervisorResponse(s)
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
