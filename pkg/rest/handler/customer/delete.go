package customer

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Delete godoc
// @Summary Delete a customer
// @Description Soft delete a customer by ID
// @Tags Customers
// @Produce json
// @Param id path int true "Customer ID"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/customers/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.WarnContext(r.Context(), "Invalid customer ID",
			slog.String("id", idStr))
		response.Error(w, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	if err := h.useCase.Delete(r.Context(), id); err != nil {
		slog.ErrorContext(r.Context(), "Failed to delete customer",
			slog.String("error", err.Error()),
			slog.Int64("customer_id", id))

		if err.Error() == "cannot delete customer with associated properties" {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
