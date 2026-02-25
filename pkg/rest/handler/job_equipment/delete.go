package job_equipment

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Delete godoc
// @Summary Delete job equipment
// @Description Delete a job equipment entry (hard delete)
// @Tags JobEquipment
// @Param jobId path int true "Job ID"
// @Param id path int true "Equipment ID"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/jobs/{jobId}/equipment/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	jobIDStr := chi.URLParam(r, "jobId")
	jobID, err := strconv.ParseInt(jobIDStr, 10, 64)
	if err != nil {
		slog.WarnContext(r.Context(), "Invalid job ID",
			slog.String("jobId", jobIDStr))
		response.Error(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.WarnContext(r.Context(), "Invalid equipment ID",
			slog.String("id", idStr))
		response.Error(w, http.StatusBadRequest, "Invalid equipment ID")
		return
	}

	if err := h.useCase.Delete(r.Context(), id, jobID); err != nil {
		if err.Error() == "job equipment not found" {
			response.Error(w, http.StatusNotFound, "Job equipment not found")
			return
		}

		if err.Error() == "equipment does not belong to this job" {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		slog.ErrorContext(r.Context(), "Failed to delete job equipment",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Failed to delete job equipment")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
