package job

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	domainJob "github.com/your-org/jvairv2/pkg/domain/job"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Update maneja la solicitud de actualización de un job
// @Summary Actualizar trabajo
// @Description Actualiza un trabajo existente. Si se cambia el technician_job_status_id y este tiene un job_status_id vinculado, el job_status_id se actualiza automáticamente
// @Tags Jobs
// @Accept json
// @Produce json
// @Param id path int true "ID del trabajo"
// @Param job body UpdateJobRequest true "Datos del trabajo"
// @Success 200 {object} JobResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/jobs/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var req UpdateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.WarnContext(r.Context(), "Invalid request body",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	// Obtener job existente para merge
	existing, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if err == domainJob.ErrJobNotFound {
			response.Error(w, http.StatusNotFound, "Trabajo no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener trabajo")
		return
	}

	// Merge: solo actualizar campos proporcionados
	j := *existing
	j.ID = id

	if req.WorkOrder != nil {
		j.WorkOrder = req.WorkOrder
	}
	if req.DateReceived != nil && *req.DateReceived != "" {
		t, err := time.Parse("01-02-2006", *req.DateReceived)
		if err != nil {
			t, err = time.Parse("2006-01-02", *req.DateReceived)
			if err != nil {
				response.Error(w, http.StatusBadRequest, "Formato de fecha de recepción inválido")
				return
			}
		}
		j.DateReceived = t
	}
	if req.JobCategoryID != nil {
		j.JobCategoryID = *req.JobCategoryID
	}
	if req.JobPriorityID != nil {
		j.JobPriorityID = *req.JobPriorityID
	}
	if req.JobStatusID != nil {
		j.JobStatusID = *req.JobStatusID
	}
	if req.TechnicianJobStatusID != nil {
		j.TechnicianJobStatusID = req.TechnicianJobStatusID
	}
	if req.WorkflowID != nil {
		j.WorkflowID = *req.WorkflowID
	}
	if req.UserID != nil {
		j.UserID = req.UserID
	}
	if req.SupervisorIDs != nil {
		j.SupervisorIDs = req.SupervisorIDs
	}
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
	if req.CompletionDate != nil && *req.CompletionDate != "" {
		t, err := time.Parse("01-02-2006", *req.CompletionDate)
		if err != nil {
			t, err = time.Parse("2006-01-02", *req.CompletionDate)
			if err != nil {
				response.Error(w, http.StatusBadRequest, "Formato de fecha de completado inválido")
				return
			}
		}
		j.CompletionDate = &t
	}
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
	if req.InstallationDueDate != nil && *req.InstallationDueDate != "" {
		t, err := time.Parse("01-02-2006", *req.InstallationDueDate)
		if err != nil {
			t, err = time.Parse("2006-01-02", *req.InstallationDueDate)
			if err != nil {
				response.Error(w, http.StatusBadRequest, "Formato de fecha de instalación inválido")
				return
			}
		}
		j.InstallationDueDate = &t
	}
	if req.WeekNumber != nil {
		j.WeekNumber = req.WeekNumber
	}
	if req.RouteNumber != nil {
		j.RouteNumber = req.RouteNumber
	}
	if req.ScheduledTimeType != nil {
		j.ScheduledTimeType = req.ScheduledTimeType
	}
	if req.ScheduledTime != nil {
		j.ScheduledTime = req.ScheduledTime
	}
	if req.InternalJobNotes != nil {
		j.InternalJobNotes = req.InternalJobNotes
	}
	if req.QuickNotes != nil {
		j.QuickNotes = req.QuickNotes
	}
	if req.JobReport != nil {
		j.JobReport = req.JobReport
	}
	if req.CageRequired != nil {
		j.CageRequired = *req.CageRequired
	}
	if req.WarrantyClaim != nil {
		j.WarrantyClaim = *req.WarrantyClaim
	}
	if req.WarrantyRegistration != nil {
		j.WarrantyRegistration = *req.WarrantyRegistration
	}
	if req.JobSalesPrice != nil {
		j.JobSalesPrice = req.JobSalesPrice
	}
	if req.MoneyTurnedIn != nil {
		j.MoneyTurnedIn = req.MoneyTurnedIn
	}
	if req.Closed != nil {
		j.Closed = *req.Closed
	}
	if req.DispatchNotes != nil {
		j.DispatchNotes = req.DispatchNotes
	}
	if req.CallLogs != nil {
		j.CallLogs = req.CallLogs
	}
	if req.CallAttempted != nil {
		j.CallAttempted = *req.CallAttempted
	}

	if err := h.useCase.Update(r.Context(), &j); err != nil {
		switch err {
		case domainJob.ErrJobNotFound:
			response.Error(w, http.StatusNotFound, "Trabajo no encontrado")
		case domainJob.ErrInvalidJobCategory,
			domainJob.ErrInvalidJobPriority,
			domainJob.ErrInvalidJobStatus,
			domainJob.ErrInvalidWorkflow,
			domainJob.ErrInvalidUser,
			domainJob.ErrInvalidTechnicianJobStatus:
			response.Error(w, http.StatusBadRequest, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "Error al actualizar trabajo")
		}
		return
	}

	// Re-fetch para obtener datos actualizados
	updated, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		response.JSON(w, http.StatusOK, toJobResponse(&j))
		return
	}

	response.JSON(w, http.StatusOK, toJobResponse(updated))
}
