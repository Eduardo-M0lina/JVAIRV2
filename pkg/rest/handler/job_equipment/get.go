package job_equipment

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Get godoc
// @Summary Get job equipment by ID
// @Description Get a job equipment entry by its ID
// @Tags job-equipment
// @Produce json
// @Param jobId path int true "Job ID"
// @Param id path int true "Equipment ID"
// @Success 200 {object} JobEquipmentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /jobs/{jobId}/equipment/{id} [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.WarnContext(r.Context(), "Invalid equipment ID",
			slog.String("id", idStr))
		response.Error(w, http.StatusBadRequest, "Invalid equipment ID")
		return
	}

	equipment, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if err.Error() == "job equipment not found" {
			response.Error(w, http.StatusNotFound, "Job equipment not found")
			return
		}

		slog.ErrorContext(r.Context(), "Failed to get job equipment",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Failed to get job equipment")
		return
	}

	response.JSON(w, http.StatusOK, toResponse(equipment))
}
