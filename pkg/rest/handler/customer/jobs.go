package customer

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// GetJobs godoc
// @Summary Get customer jobs
// @Description Get all jobs for a specific customer
// @Tags Customers
// @Produce json
// @Param id path int true "Customer ID"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Param search query string false "Search term"
// @Param closed query string false "Filter by closed status: 0 (open), 1 (closed), all" default(0)
// @Param job_category_id query int false "Filter by job category ID"
// @Param job_status_id query int false "Filter by job status ID"
// @Param job_priority_id query int false "Filter by job priority ID"
// @Success 200 {object} response.PaginatedResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/customers/{id}/jobs [get]
// @Security BearerAuth
func (h *Handler) GetJobs(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	customerID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.WarnContext(r.Context(), "Invalid customer ID",
			slog.String("id", idStr))
		response.Error(w, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	// Verificar que el customer existe
	_, err = h.useCase.GetByID(r.Context(), customerID)
	if err != nil {
		if err.Error() == "not found" {
			response.Error(w, http.StatusNotFound, "Customer not found")
			return
		}
		slog.ErrorContext(r.Context(), "Failed to get customer",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Failed to get customer")
		return
	}

	// Parsear parámetros de paginación
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize < 1 {
		pageSize = 10
	}

	// Crear filtros con customerId
	filters := map[string]interface{}{
		"customer_id": customerID,
	}

	// Agregar búsqueda si existe
	if search := r.URL.Query().Get("search"); search != "" {
		filters["search"] = search
	}

	// Filtro por closed
	if closed := r.URL.Query().Get("closed"); closed != "" {
		filters["closed"] = closed
	}

	// Filtro por categoría
	if jobCategoryIDStr := r.URL.Query().Get("job_category_id"); jobCategoryIDStr != "" {
		if jobCategoryID, err := strconv.ParseInt(jobCategoryIDStr, 10, 64); err == nil {
			filters["job_category_id"] = jobCategoryID
		}
	}

	// Filtro por estado
	if jobStatusIDStr := r.URL.Query().Get("job_status_id"); jobStatusIDStr != "" {
		if jobStatusID, err := strconv.ParseInt(jobStatusIDStr, 10, 64); err == nil {
			filters["job_status_id"] = jobStatusID
		}
	}

	// Filtro por prioridad
	if jobPriorityIDStr := r.URL.Query().Get("job_priority_id"); jobPriorityIDStr != "" {
		if jobPriorityID, err := strconv.ParseInt(jobPriorityIDStr, 10, 64); err == nil {
			filters["job_priority_id"] = jobPriorityID
		}
	}

	// Obtener jobs usando el use case de job
	jobs, total, err := h.jobUC.List(r.Context(), filters, page, pageSize)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to list customer jobs",
			slog.String("error", err.Error()),
			slog.Int64("customerId", customerID))
		response.Error(w, http.StatusInternalServerError, "Failed to list jobs")
		return
	}

	// Convertir a response
	items := make([]map[string]interface{}, len(jobs))
	for i, j := range jobs {
		items[i] = map[string]interface{}{
			"id":                    j.ID,
			"workOrder":             j.WorkOrder,
			"dateReceived":          j.DateReceived.Format("2006-01-02T15:04:05Z07:00"),
			"jobCategoryId":         j.JobCategoryID,
			"jobPriorityId":         j.JobPriorityID,
			"jobStatusId":           j.JobStatusID,
			"technicianJobStatusId": j.TechnicianJobStatusID,
			"workflowId":            j.WorkflowID,
			"propertyId":            j.PropertyID,
			"userId":                j.UserID,
			"supervisorIds":         j.SupervisorIDs,
			"weekNumber":            j.WeekNumber,
			"routeNumber":           j.RouteNumber,
			"scheduledTimeType":     j.ScheduledTimeType,
			"scheduledTime":         j.ScheduledTime,
			"internalJobNotes":      j.InternalJobNotes,
			"quickNotes":            j.QuickNotes,
			"jobReport":             j.JobReport,
			"cageRequired":          j.CageRequired,
			"warrantyClaim":         j.WarrantyClaim,
			"warrantyRegistration":  j.WarrantyRegistration,
			"jobSalesPrice":         j.JobSalesPrice,
			"moneyTurnedIn":         j.MoneyTurnedIn,
			"closed":                j.Closed,
			"dispatchNotes":         j.DispatchNotes,
			"callLogs":              j.CallLogs,
			"callAttempted":         j.CallAttempted,
		}
		if j.DispatchDate != nil {
			items[i]["dispatchDate"] = j.DispatchDate.Format("2006-01-02T15:04:05Z07:00")
		}
		if j.CompletionDate != nil {
			items[i]["completionDate"] = j.CompletionDate.Format("2006-01-02T15:04:05Z07:00")
		}
		if j.InstallationDueDate != nil {
			items[i]["installationDueDate"] = j.InstallationDueDate.Format("2006-01-02T15:04:05Z07:00")
		}
		if j.DueDate != nil {
			items[i]["dueDate"] = j.DueDate.Format("2006-01-02T15:04:05Z07:00")
		}
		if j.CreatedAt != nil {
			items[i]["createdAt"] = j.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
		}
		if j.UpdatedAt != nil {
			items[i]["updatedAt"] = j.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
		}
	}

	totalPages := (total + pageSize - 1) / pageSize

	response.JSON(w, http.StatusOK, response.PaginatedResponse{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	})
}
