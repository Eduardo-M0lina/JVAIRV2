package supervisor

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Delete godoc
// @Summary Delete a supervisor
// @Description Soft delete a supervisor by ID
// @Tags supervisors
// @Produce json
// @Param id path int true "Supervisor ID"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /supervisors/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.WarnContext(r.Context(), "Invalid supervisor ID",
			slog.String("id", idStr))
		response.Error(w, http.StatusBadRequest, "Invalid supervisor ID")
		return
	}

	if err := h.useCase.Delete(r.Context(), id); err != nil {
		slog.ErrorContext(r.Context(), "Failed to delete supervisor",
			slog.String("error", err.Error()),
			slog.Int64("supervisor_id", id))

		if err.Error() == "supervisor already deleted" {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
