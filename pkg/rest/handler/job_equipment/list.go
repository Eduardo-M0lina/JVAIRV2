package job_equipment

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// List godoc
// @Summary List job equipment
// @Description Get a list of equipment for a job with optional type filter
// @Tags job-equipment
// @Produce json
// @Param jobId path int true "Job ID"
// @Param type query string false "Filter by type (current, new)"
// @Success 200 {array} JobEquipmentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /jobs/{jobId}/equipment [get]
// @Security BearerAuth
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	jobIDStr := chi.URLParam(r, "jobId")
	jobID, err := strconv.ParseInt(jobIDStr, 10, 64)
	if err != nil {
		slog.WarnContext(r.Context(), "Invalid job ID",
			slog.String("jobId", jobIDStr))
		response.Error(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	equipmentType := r.URL.Query().Get("type")

	equipment, err := h.useCase.List(r.Context(), jobID, equipmentType)
	if err != nil {
		if err.Error() == "invalid job_id" ||
			err.Error() == "type must be one of: current, new" {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		slog.ErrorContext(r.Context(), "Failed to list job equipment",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Failed to list job equipment")
		return
	}

	items := make([]JobEquipmentResponse, len(equipment))
	for i, e := range equipment {
		items[i] = toResponse(e)
	}

	response.JSON(w, http.StatusOK, items)
}
