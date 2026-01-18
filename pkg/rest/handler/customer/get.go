package customer

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Get godoc
// @Summary Get customer by ID
// @Description Get a customer by its ID
// @Tags customers
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200 {object} CustomerResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /customers/{id} [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.WarnContext(r.Context(), "Invalid customer ID",
			slog.String("id", idStr))
		response.Error(w, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	c, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to get customer",
			slog.String("error", err.Error()),
			slog.Int64("customer_id", id))
		response.Error(w, http.StatusNotFound, "Customer not found")
		return
	}

	response.JSON(w, http.StatusOK, toCustomerResponse(c))
}
