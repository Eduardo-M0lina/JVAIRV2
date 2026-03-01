package property

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Delete godoc
// @Summary Delete property
// @Description Delete a property (soft delete)
// @Tags Properties
// @Param id path int true "Property ID"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/properties/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.WarnContext(r.Context(), "Invalid property ID",
			slog.String("id", idStr))
		response.Error(w, http.StatusBadRequest, "Invalid property ID")
		return
	}

	if err := h.useCase.Delete(r.Context(), id); err != nil {
		if err.Error() == "property not found" {
			response.Error(w, http.StatusNotFound, "Property not found")
			return
		}

		if err.Error() == "property already deleted" ||
			err.Error() == "cannot delete property with associated jobs" {
			slog.WarnContext(r.Context(), "Cannot delete property",
				slog.String("error", err.Error()))
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		slog.ErrorContext(r.Context(), "Failed to delete property",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Failed to delete property")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
