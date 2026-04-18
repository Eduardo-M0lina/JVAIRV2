package job

import (
	"fmt"
	"strings"
	"time"
)

// JobCategory representa la categoría de un job
type JobCategory struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
	Type  string `json:"type"`
}

// JobStatus representa el estado de un job
type JobStatus struct {
	ID    int64   `json:"id"`
	Label string  `json:"label"`
	Class *string `json:"class,omitempty"`
}

// JobPriority representa la prioridad de un job
type JobPriority struct {
	ID    int64   `json:"id"`
	Label string  `json:"label"`
	Order int     `json:"order"`
	Class *string `json:"class,omitempty"`
}

// Job representa un trabajo en el sistema
type Job struct {
	ID                    int64        `json:"id"`
	WorkOrder             *string      `json:"workOrder,omitempty"`
	DateReceived          time.Time    `json:"dateReceived"`
	JobCategoryID         int64        `json:"jobCategoryId"`
	Category              *JobCategory `json:"category,omitempty"`
	JobPriorityID         int64        `json:"jobPriorityId"`
	Priority              *JobPriority `json:"priority,omitempty"`
	JobStatusID           int64        `json:"jobStatusId"`
	Status                *JobStatus   `json:"status,omitempty"`
	TechnicianJobStatusID *int64       `json:"technicianJobStatusId,omitempty"`
	WorkflowID            int64        `json:"workflowId"`
	PropertyID            int64        `json:"propertyId"`
	UserID                *int64       `json:"userId,omitempty"`
	SupervisorIDs         *string      `json:"supervisorIds,omitempty"`
	DispatchDate          *time.Time   `json:"dispatchDate,omitempty"`
	CompletionDate        *time.Time   `json:"completionDate,omitempty"`
	WeekNumber            *int         `json:"weekNumber,omitempty"`
	RouteNumber           *int         `json:"routeNumber,omitempty"`
	ScheduledTimeType     *string      `json:"scheduledTimeType,omitempty"`
	ScheduledTime         *string      `json:"scheduledTime,omitempty"`
	InternalJobNotes      *string      `json:"internalJobNotes,omitempty"`
	QuickNotes            *string      `json:"quickNotes,omitempty"`
	JobReport             *string      `json:"jobReport,omitempty"`
	InstallationDueDate   *time.Time   `json:"installationDueDate,omitempty"`
	CageRequired          bool         `json:"cageRequired"`
	WarrantyClaim         bool         `json:"warrantyClaim"`
	WarrantyRegistration  bool         `json:"warrantyRegistration"`
	JobSalesPrice         *float64     `json:"jobSalesPrice,omitempty"`
	MoneyTurnedIn         *float64     `json:"moneyTurnedIn,omitempty"`
	Closed                bool         `json:"closed"`
	DispatchNotes         *string      `json:"dispatchNotes,omitempty"`
	CallLogs              *string      `json:"callLogs,omitempty"`
	DueDate               *time.Time   `json:"dueDate,omitempty"`
	CallAttempted         bool         `json:"callAttempted"`
	CreatedAt             *time.Time   `json:"createdAt,omitempty"`
	UpdatedAt             *time.Time   `json:"updatedAt,omitempty"`
	DeletedAt             *time.Time   `json:"deletedAt,omitempty"`
}

// ValidateCreate valida los campos requeridos para crear un job
func (j *Job) ValidateCreate() error {
	if j.JobCategoryID == 0 {
		return fmt.Errorf("job_category_id is required")
	}

	if j.JobPriorityID == 0 {
		return fmt.Errorf("job_priority_id is required")
	}

	if j.PropertyID == 0 {
		return fmt.Errorf("property_id is required")
	}

	if j.WorkOrder != nil && strings.TrimSpace(*j.WorkOrder) == "" {
		j.WorkOrder = nil
	}

	return nil
}

// ValidateUpdate valida los campos para actualizar un job
func (j *Job) ValidateUpdate() error {
	if j.ID == 0 {
		return fmt.Errorf("id is required")
	}

	return nil
}

// IsDeleted verifica si el job está eliminado
func (j *Job) IsDeleted() bool {
	return j.DeletedAt != nil
}

// IsClosed verifica si el job está cerrado
func (j *Job) IsClosed() bool {
	return j.Closed
}
