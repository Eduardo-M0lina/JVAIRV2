package supervisor

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Get godoc
// @Summary Get supervisor by ID
// @Description Get a supervisor by its ID
// @Tags Supervisors
// @Produce json
// @Param id path int true "Supervisor ID"
// @Success 200 {object} SupervisorResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/supervisors/{id} [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.WarnContext(r.Context(), "Invalid supervisor ID",
			slog.String("id", idStr))
		response.Error(w, http.StatusBadRequest, "Invalid supervisor ID")
		return
	}

	s, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to get supervisor",
			slog.String("error", err.Error()),
			slog.Int64("supervisorId", id))
		response.Error(w, http.StatusNotFound, "Supervisor not found")
		return
	}

	response.JSON(w, http.StatusOK, toSupervisorResponse(s))
}
