package job

import (
	"context"
	"net/http"
	"strconv"

	domainJob "github.com/angumol/jvairv2/pkg/domain/job"
	domainJobActivityLog "github.com/angumol/jvairv2/pkg/domain/job_activity_log"
	domainJobSMS "github.com/angumol/jvairv2/pkg/domain/job_sms"
	domainUser "github.com/angumol/jvairv2/pkg/domain/user"
	"github.com/angumol/jvairv2/pkg/rest/middleware"
	"github.com/go-chi/chi/v5"
)

// Handler maneja las peticiones HTTP para jobs
type Handler struct {
	useCase            domainJob.Service
	activityLogService domainJobActivityLog.Service
	emailService       EmailService
	smsSender          SMSSender
	jobSMSService      domainJobSMS.Service
}

// EmailService define la interfaz para envío de emails
type EmailService interface {
	SendDispatchEmail(ctx context.Context, jobID int64, recipients []string) error
	SendDispatchSupervisorEmail(ctx context.Context, jobID int64, subject string, body string, recipients []string) error
}

// SMSSender define la interfaz para envío de SMS
type SMSSender interface {
	SendSMS(ctx context.Context, to string, body string) error
}

// NewHandler crea una nueva instancia del handler de jobs
func NewHandler(useCase domainJob.Service, activityLogService domainJobActivityLog.Service, emailService EmailService, smsSender SMSSender, jobSMSService domainJobSMS.Service) *Handler {
	return &Handler{
		useCase:            useCase,
		activityLogService: activityLogService,
		emailService:       emailService,
		smsSender:          smsSender,
		jobSMSService:      jobSMSService,
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
		r.Post("/{id}/dispatch-email", h.SendDispatchEmail)
		r.Post("/{id}/dispatch-supervisor-email", h.SendDispatchSupervisorEmail)
		r.Post("/{id}/dispatch-sms", h.SendDispatchSMS)
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

// logActivity registra automáticamente una actividad en el job
func (h *Handler) logActivity(ctx context.Context, jobID int64, activityType, message string) {
	// Obtener usuario del contexto
	userFromCtx, ok := ctx.Value(middleware.UserContextKey).(*domainUser.User)
	if !ok || userFromCtx == nil {
		// Si no hay usuario en contexto, no registrar actividad
		return
	}

	activity := &domainJobActivityLog.JobActivityLog{
		JobID:  jobID,
		UserID: userFromCtx.ID,
		Type:   activityType,
		Log:    message,
	}

	// Intentar crear la actividad, pero no fallar si hay error
	_ = h.activityLogService.Create(ctx, activity)
}
