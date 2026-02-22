package job

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	domainJob "github.com/your-org/jvairv2/pkg/domain/job"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Create maneja la solicitud de creación de un job
// @Summary Crear trabajo
// @Description Crea un nuevo trabajo. El workflow y status inicial se asignan automáticamente desde el customer de la propiedad
// @Tags Jobs
// @Accept json
// @Produce json
// @Param job body CreateJobRequest true "Datos del trabajo"
// @Success 201 {object} JobResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/jobs [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.WarnContext(r.Context(), "Invalid request body",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	j := &domainJob.Job{
		WorkOrder:         req.WorkOrder,
		JobCategoryID:     req.JobCategoryID,
		JobPriorityID:     req.JobPriorityID,
		PropertyID:        req.PropertyID,
		UserID:            req.UserID,
		SupervisorIDs:     req.SupervisorIDs,
		WeekNumber:        req.WeekNumber,
		RouteNumber:       req.RouteNumber,
		ScheduledTimeType: req.ScheduledTimeType,
		ScheduledTime:     req.ScheduledTime,
		DispatchNotes:     req.DispatchNotes,
		QuickNotes:        req.QuickNotes,
		InternalJobNotes:  req.InternalJobNotes,
		JobSalesPrice:     req.JobSalesPrice,
	}

	// Parsear fecha de recepción
	if req.DateReceived != nil && *req.DateReceived != "" {
		t, err := time.Parse("01-02-2006", *req.DateReceived)
		if err != nil {
			// Intentar formato ISO
			t, err = time.Parse("2006-01-02", *req.DateReceived)
			if err != nil {
				response.Error(w, http.StatusBadRequest, "Formato de fecha de recepción inválido (usar MM-DD-YYYY o YYYY-MM-DD)")
				return
			}
		}
		j.DateReceived = t
	}

	// Parsear fecha de dispatch
	if req.DispatchDate != nil && *req.DispatchDate != "" {
		t, err := time.Parse("01-02-2006", *req.DispatchDate)
		if err != nil {
			t, err = time.Parse("2006-01-02", *req.DispatchDate)
			if err != nil {
				response.Error(w, http.StatusBadRequest, "Formato de fecha de dispatch inválido")
				return
			}
		}
		j.DispatchDate = &t
	}

	// Parsear due date
	if req.DueDate != nil && *req.DueDate != "" {
		t, err := time.Parse("01-02-2006", *req.DueDate)
		if err != nil {
			t, err = time.Parse("2006-01-02", *req.DueDate)
			if err != nil {
				response.Error(w, http.StatusBadRequest, "Formato de due date inválido")
				return
			}
		}
		j.DueDate = &t
	}

	// Booleans opcionales
	if req.CageRequired != nil {
		j.CageRequired = *req.CageRequired
	}
	if req.WarrantyClaim != nil {
		j.WarrantyClaim = *req.WarrantyClaim
	}
	if req.WarrantyRegistration != nil {
		j.WarrantyRegistration = *req.WarrantyRegistration
	}

	if err := h.useCase.Create(r.Context(), j); err != nil {
		switch err {
		case domainJob.ErrInvalidJobCategory,
			domainJob.ErrInvalidJobPriority,
			domainJob.ErrInvalidProperty,
			domainJob.ErrInvalidUser,
			domainJob.ErrInvalidWorkflow,
			domainJob.ErrInvalidJobStatus,
			domainJob.ErrWorkflowHasNoStatuses:
			response.Error(w, http.StatusBadRequest, err.Error())
		default:
			if err.Error() == "job_category_id is required" ||
				err.Error() == "job_priority_id is required" ||
				err.Error() == "property_id is required" {
				response.Error(w, http.StatusBadRequest, err.Error())
			} else {
				response.Error(w, http.StatusInternalServerError, "Error al crear trabajo")
			}
		}
		return
	}

	response.JSON(w, http.StatusCreated, toJobResponse(j))
}
