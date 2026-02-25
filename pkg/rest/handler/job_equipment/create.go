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

// Create godoc
// @Summary Create job equipment
// @Description Create a new equipment entry for a job
// @Tags JobEquipment
// @Accept json
// @Produce json
// @Param jobId path int true "Job ID"
// @Param equipment body CreateJobEquipmentRequest true "Equipment data"
// @Success 201 {object} JobEquipmentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/jobs/{jobId}/equipment [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	jobIDStr := chi.URLParam(r, "jobId")
	jobID, err := strconv.ParseInt(jobIDStr, 10, 64)
	if err != nil {
		slog.WarnContext(r.Context(), "Invalid job ID",
			slog.String("jobId", jobIDStr))
		response.Error(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	var req CreateJobEquipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.WarnContext(r.Context(), "Invalid request body",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	eq := &domain.JobEquipment{
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

	if err := h.useCase.Create(r.Context(), eq); err != nil {
		if err.Error() == "invalid job_id" {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		slog.ErrorContext(r.Context(), "Failed to create job equipment",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Failed to create job equipment")
		return
	}

	response.JSON(w, http.StatusCreated, toResponse(eq))
}
