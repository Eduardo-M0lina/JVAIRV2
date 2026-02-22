package job

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	domainJob "github.com/your-org/jvairv2/pkg/domain/job"
)

// Handler maneja las peticiones HTTP para jobs
type Handler struct {
	useCase domainJob.Service
}

// NewHandler crea una nueva instancia del handler de jobs
func NewHandler(useCase domainJob.Service) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

// RegisterRoutes registra las rutas del handler
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/jobs", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
		r.Put("/{id}/close", h.Close)
	})
}

// CreateJobRequest representa la solicitud para crear un job
type CreateJobRequest struct {
	WorkOrder            *string  `json:"workOrder,omitempty"`
	DateReceived         *string  `json:"dateReceived,omitempty"`
	JobCategoryID        int64    `json:"jobCategoryId"`
	JobPriorityID        int64    `json:"jobPriorityId"`
	PropertyID           int64    `json:"propertyId"`
	UserID               *int64   `json:"userId,omitempty"`
	SupervisorIDs        *string  `json:"supervisorIds,omitempty"`
	DispatchDate         *string  `json:"dispatchDate,omitempty"`
	DueDate              *string  `json:"dueDate,omitempty"`
	WeekNumber           *int     `json:"weekNumber,omitempty"`
	RouteNumber          *int     `json:"routeNumber,omitempty"`
	ScheduledTimeType    *string  `json:"scheduledTimeType,omitempty"`
	ScheduledTime        *string  `json:"scheduledTime,omitempty"`
	DispatchNotes        *string  `json:"dispatchNotes,omitempty"`
	QuickNotes           *string  `json:"quickNotes,omitempty"`
	InternalJobNotes     *string  `json:"internalJobNotes,omitempty"`
	JobSalesPrice        *float64 `json:"jobSalesPrice,omitempty"`
	CageRequired         *bool    `json:"cageRequired,omitempty"`
	WarrantyClaim        *bool    `json:"warrantyClaim,omitempty"`
	WarrantyRegistration *bool    `json:"warrantyRegistration,omitempty"`
}

// UpdateJobRequest representa la solicitud para actualizar un job
type UpdateJobRequest struct {
	WorkOrder             *string  `json:"workOrder,omitempty"`
	DateReceived          *string  `json:"dateReceived,omitempty"`
	JobCategoryID         *int64   `json:"jobCategoryId,omitempty"`
	JobPriorityID         *int64   `json:"jobPriorityId,omitempty"`
	JobStatusID           *int64   `json:"jobStatusId,omitempty"`
	TechnicianJobStatusID *int64   `json:"technicianJobStatusId,omitempty"`
	WorkflowID            *int64   `json:"workflowId,omitempty"`
	UserID                *int64   `json:"userId,omitempty"`
	SupervisorIDs         *string  `json:"supervisorIds,omitempty"`
	DispatchDate          *string  `json:"dispatchDate,omitempty"`
	CompletionDate        *string  `json:"completionDate,omitempty"`
	DueDate               *string  `json:"dueDate,omitempty"`
	WeekNumber            *int     `json:"weekNumber,omitempty"`
	RouteNumber           *int     `json:"routeNumber,omitempty"`
	ScheduledTimeType     *string  `json:"scheduledTimeType,omitempty"`
	ScheduledTime         *string  `json:"scheduledTime,omitempty"`
	InternalJobNotes      *string  `json:"internalJobNotes,omitempty"`
	QuickNotes            *string  `json:"quickNotes,omitempty"`
	JobReport             *string  `json:"jobReport,omitempty"`
	InstallationDueDate   *string  `json:"installationDueDate,omitempty"`
	CageRequired          *bool    `json:"cageRequired,omitempty"`
	WarrantyClaim         *bool    `json:"warrantyClaim,omitempty"`
	WarrantyRegistration  *bool    `json:"warrantyRegistration,omitempty"`
	JobSalesPrice         *float64 `json:"jobSalesPrice,omitempty"`
	MoneyTurnedIn         *float64 `json:"moneyTurnedIn,omitempty"`
	Closed                *bool    `json:"closed,omitempty"`
	DispatchNotes         *string  `json:"dispatchNotes,omitempty"`
	CallLogs              *string  `json:"callLogs,omitempty"`
	CallAttempted         *bool    `json:"callAttempted,omitempty"`
}

// CloseJobRequest representa la solicitud para cerrar un job
type CloseJobRequest struct {
	JobStatusID int64 `json:"jobStatusId"`
}

// JobResponse representa la respuesta de un job
type JobResponse struct {
	ID                    int64    `json:"id"`
	WorkOrder             *string  `json:"workOrder,omitempty"`
	DateReceived          string   `json:"dateReceived"`
	JobCategoryID         int64    `json:"jobCategoryId"`
	JobPriorityID         int64    `json:"jobPriorityId"`
	JobStatusID           int64    `json:"jobStatusId"`
	TechnicianJobStatusID *int64   `json:"technicianJobStatusId,omitempty"`
	WorkflowID            int64    `json:"workflowId"`
	PropertyID            int64    `json:"propertyId"`
	UserID                *int64   `json:"userId,omitempty"`
	SupervisorIDs         *string  `json:"supervisorIds,omitempty"`
	DispatchDate          *string  `json:"dispatchDate,omitempty"`
	CompletionDate        *string  `json:"completionDate,omitempty"`
	WeekNumber            *int     `json:"weekNumber,omitempty"`
	RouteNumber           *int     `json:"routeNumber,omitempty"`
	ScheduledTimeType     *string  `json:"scheduledTimeType,omitempty"`
	ScheduledTime         *string  `json:"scheduledTime,omitempty"`
	InternalJobNotes      *string  `json:"internalJobNotes,omitempty"`
	QuickNotes            *string  `json:"quickNotes,omitempty"`
	JobReport             *string  `json:"jobReport,omitempty"`
	InstallationDueDate   *string  `json:"installationDueDate,omitempty"`
	CageRequired          bool     `json:"cageRequired"`
	WarrantyClaim         bool     `json:"warrantyClaim"`
	WarrantyRegistration  bool     `json:"warrantyRegistration"`
	JobSalesPrice         *float64 `json:"jobSalesPrice,omitempty"`
	MoneyTurnedIn         *float64 `json:"moneyTurnedIn,omitempty"`
	Closed                bool     `json:"closed"`
	DispatchNotes         *string  `json:"dispatchNotes,omitempty"`
	CallLogs              *string  `json:"callLogs,omitempty"`
	DueDate               *string  `json:"dueDate,omitempty"`
	CallAttempted         bool     `json:"callAttempted"`
	CreatedAt             string   `json:"createdAt,omitempty"`
	UpdatedAt             string   `json:"updatedAt,omitempty"`
}

const timeFormat = "2006-01-02T15:04:05Z07:00"

func toJobResponse(j *domainJob.Job) JobResponse {
	resp := JobResponse{
		ID:                    j.ID,
		WorkOrder:             j.WorkOrder,
		DateReceived:          j.DateReceived.Format(timeFormat),
		JobCategoryID:         j.JobCategoryID,
		JobPriorityID:         j.JobPriorityID,
		JobStatusID:           j.JobStatusID,
		TechnicianJobStatusID: j.TechnicianJobStatusID,
		WorkflowID:            j.WorkflowID,
		PropertyID:            j.PropertyID,
		UserID:                j.UserID,
		SupervisorIDs:         j.SupervisorIDs,
		WeekNumber:            j.WeekNumber,
		RouteNumber:           j.RouteNumber,
		ScheduledTimeType:     j.ScheduledTimeType,
		ScheduledTime:         j.ScheduledTime,
		InternalJobNotes:      j.InternalJobNotes,
		QuickNotes:            j.QuickNotes,
		JobReport:             j.JobReport,
		CageRequired:          j.CageRequired,
		WarrantyClaim:         j.WarrantyClaim,
		WarrantyRegistration:  j.WarrantyRegistration,
		JobSalesPrice:         j.JobSalesPrice,
		MoneyTurnedIn:         j.MoneyTurnedIn,
		Closed:                j.Closed,
		DispatchNotes:         j.DispatchNotes,
		CallLogs:              j.CallLogs,
		CallAttempted:         j.CallAttempted,
	}

	if j.DispatchDate != nil {
		s := j.DispatchDate.Format(timeFormat)
		resp.DispatchDate = &s
	}
	if j.CompletionDate != nil {
		s := j.CompletionDate.Format(timeFormat)
		resp.CompletionDate = &s
	}
	if j.InstallationDueDate != nil {
		s := j.InstallationDueDate.Format(timeFormat)
		resp.InstallationDueDate = &s
	}
	if j.DueDate != nil {
		s := j.DueDate.Format(timeFormat)
		resp.DueDate = &s
	}
	if j.CreatedAt != nil {
		resp.CreatedAt = j.CreatedAt.Format(timeFormat)
	}
	if j.UpdatedAt != nil {
		resp.UpdatedAt = j.UpdatedAt.Format(timeFormat)
	}

	return resp
}

func parseFilters(r *http.Request) map[string]interface{} {
	filters := make(map[string]interface{})

	if search := r.URL.Query().Get("search"); search != "" {
		filters["search"] = search
	}

	if closed := r.URL.Query().Get("closed"); closed != "" {
		filters["closed"] = closed
	}

	if jobCategoryIDStr := r.URL.Query().Get("jobCategoryId"); jobCategoryIDStr != "" {
		if id, err := strconv.ParseInt(jobCategoryIDStr, 10, 64); err == nil {
			filters["job_category_id"] = id
		}
	}

	if jobStatusIDStr := r.URL.Query().Get("jobStatusId"); jobStatusIDStr != "" {
		if id, err := strconv.ParseInt(jobStatusIDStr, 10, 64); err == nil {
			filters["job_status_id"] = id
		}
	}

	if jobPriorityIDStr := r.URL.Query().Get("jobPriorityId"); jobPriorityIDStr != "" {
		if id, err := strconv.ParseInt(jobPriorityIDStr, 10, 64); err == nil {
			filters["job_priority_id"] = id
		}
	}

	if userID := r.URL.Query().Get("userId"); userID != "" {
		if userID == "unassigned" {
			filters["user_id"] = "unassigned"
		} else if id, err := strconv.ParseInt(userID, 10, 64); err == nil {
			filters["user_id"] = id
		}
	}

	if propertyIDStr := r.URL.Query().Get("propertyId"); propertyIDStr != "" {
		if id, err := strconv.ParseInt(propertyIDStr, 10, 64); err == nil {
			filters["property_id"] = id
		}
	}

	if workflowIDStr := r.URL.Query().Get("workflowId"); workflowIDStr != "" {
		if id, err := strconv.ParseInt(workflowIDStr, 10, 64); err == nil {
			filters["workflow_id"] = id
		}
	}

	if sort := r.URL.Query().Get("sort"); sort != "" {
		filters["sort"] = sort
	}

	if direction := r.URL.Query().Get("direction"); direction != "" {
		filters["direction"] = direction
	}

	return filters
}
