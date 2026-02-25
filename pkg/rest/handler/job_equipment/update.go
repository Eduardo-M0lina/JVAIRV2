package job_equipment

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	domain "github.com/your-org/jvairv2/pkg/domain/job_equipment"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Update godoc
// @Summary Update job equipment
// @Description Update an existing job equipment entry
// @Tags JobEquipment
// @Accept json
// @Produce json
// @Param jobId path int true "Job ID"
// @Param id path int true "Equipment ID"
// @Param equipment body UpdateJobEquipmentRequest true "Equipment data"
// @Success 200 {object} JobEquipmentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/jobs/{jobId}/equipment/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
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

	var req UpdateJobEquipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.WarnContext(r.Context(), "Invalid request body",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	eq := &domain.JobEquipment{
		ID:                  id,
		JobID:               jobID,
		Type:                req.Type,
		Area:                req.Area,
		OutdoorBrand:        req.OutdoorBrand,
		OutdoorModel:        req.OutdoorModel,
		OutdoorSerial:       req.OutdoorSerial,
		OutdoorInstalled:    parseTimePtr(req.OutdoorInstalled),
		FurnaceBrand:        req.FurnaceBrand,
		FurnaceModel:        req.FurnaceModel,
		FurnaceSerial:       req.FurnaceSerial,
		FurnaceInstalled:    parseTimePtr(req.FurnaceInstalled),
		EvaporatorBrand:     req.EvaporatorBrand,
		EvaporatorModel:     req.EvaporatorModel,
		EvaporatorSerial:    req.EvaporatorSerial,
		EvaporatorInstalled: parseTimePtr(req.EvaporatorInstalled),
		AirHandlerBrand:     req.AirHandlerBrand,
		AirHandlerModel:     req.AirHandlerModel,
		AirHandlerSerial:    req.AirHandlerSerial,
		AirHandlerInstalled: parseTimePtr(req.AirHandlerInstalled),
	}

	if err := eq.Validate(); err != nil {
		slog.WarnContext(r.Context(), "Validation failed",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.useCase.Update(r.Context(), eq); err != nil {
		if err.Error() == "job equipment not found" {
			response.Error(w, http.StatusNotFound, "Job equipment not found")
			return
		}

		if err.Error() == "invalid job_id" ||
			err.Error() == "equipment does not belong to this job" {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		slog.ErrorContext(r.Context(), "Failed to update job equipment",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Failed to update job equipment")
		return
	}

	updatedEq, _ := h.useCase.GetByID(r.Context(), id)
	response.JSON(w, http.StatusOK, toResponse(updatedEq))
}
