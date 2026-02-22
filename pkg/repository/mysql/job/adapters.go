package job

import (
	"context"
	"database/sql"
	"log/slog"

	domainJob "github.com/your-org/jvairv2/pkg/domain/job"
)

// JobCategoryCheckerAdapter adapta el repositorio de job_category para el checker del job use case
type JobCategoryCheckerAdapter struct {
	db *sql.DB
}

func NewJobCategoryCheckerAdapter(db *sql.DB) domainJob.JobCategoryChecker {
	return &JobCategoryCheckerAdapter{db: db}
}

func (a *JobCategoryCheckerAdapter) GetByID(ctx context.Context, id int64) (interface{}, error) {
	var exists bool
	err := a.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM job_categories WHERE id = ? AND is_active = 1)", id).Scan(&exists)
	if err != nil || !exists {
		return nil, domainJob.ErrInvalidJobCategory
	}
	return true, nil
}

// JobPriorityCheckerAdapter adapta el repositorio de job_priority para el checker del job use case
type JobPriorityCheckerAdapter struct {
	db *sql.DB
}

func NewJobPriorityCheckerAdapter(db *sql.DB) domainJob.JobPriorityChecker {
	return &JobPriorityCheckerAdapter{db: db}
}

func (a *JobPriorityCheckerAdapter) GetByID(ctx context.Context, id int64) (interface{}, error) {
	var exists bool
	err := a.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM job_priorities WHERE id = ? AND is_active = 1)", id).Scan(&exists)
	if err != nil || !exists {
		return nil, domainJob.ErrInvalidJobPriority
	}
	return true, nil
}

// JobStatusCheckerAdapter adapta el repositorio de job_status para el checker del job use case
type JobStatusCheckerAdapter struct {
	db *sql.DB
}

func NewJobStatusCheckerAdapter(db *sql.DB) domainJob.JobStatusChecker {
	return &JobStatusCheckerAdapter{db: db}
}

func (a *JobStatusCheckerAdapter) GetByID(ctx context.Context, id int64) (interface{}, error) {
	var exists bool
	err := a.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM job_statuses WHERE id = ?)", id).Scan(&exists)
	if err != nil || !exists {
		return nil, domainJob.ErrInvalidJobStatus
	}
	return true, nil
}

// WorkflowCheckerAdapter adapta el repositorio de workflow para el checker del job use case
type WorkflowCheckerAdapter struct {
	db *sql.DB
}

func NewWorkflowCheckerAdapter(db *sql.DB) domainJob.WorkflowChecker {
	return &WorkflowCheckerAdapter{db: db}
}

func (a *WorkflowCheckerAdapter) GetByID(ctx context.Context, id int64) (interface{}, error) {
	var exists bool
	err := a.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM workflows WHERE id = ? AND is_active = 1)", id).Scan(&exists)
	if err != nil || !exists {
		return nil, domainJob.ErrInvalidWorkflow
	}
	return true, nil
}

// GetInitialStatusID obtiene el primer status del workflow (ordenado por order)
func (a *WorkflowCheckerAdapter) GetInitialStatusID(ctx context.Context, workflowID int64) (int64, error) {
	// Primero verificar cuántos statuses tiene el workflow
	var count int
	_ = a.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM job_status_workflow WHERE workflow_id = ?",
		workflowID,
	).Scan(&count)
	slog.InfoContext(ctx, "Checking workflow statuses",
		slog.Int64("workflowId", workflowID),
		slog.Int("statusCount", count))

	var statusID int64
	err := a.db.QueryRowContext(ctx,
		"SELECT job_status_id FROM job_status_workflow WHERE workflow_id = ? ORDER BY `order` ASC LIMIT 1",
		workflowID,
	).Scan(&statusID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get initial status for workflow",
			slog.Int64("workflowId", workflowID),
			slog.Int("statusCount", count),
			slog.String("error", err.Error()))
		if count == 0 {
			return 0, domainJob.ErrWorkflowHasNoStatuses
		}
		return 0, domainJob.ErrInvalidJobStatus
	}
	slog.InfoContext(ctx, "Resolved initial status for workflow",
		slog.Int64("workflowId", workflowID),
		slog.Int64("initialStatusId", statusID))
	return statusID, nil
}

// PropertyCheckerAdapter adapta el repositorio de property para el checker del job use case
type PropertyCheckerAdapter struct {
	db *sql.DB
}

func NewPropertyCheckerAdapter(db *sql.DB) domainJob.PropertyChecker {
	return &PropertyCheckerAdapter{db: db}
}

func (a *PropertyCheckerAdapter) GetByID(ctx context.Context, id int64) (interface{}, error) {
	var exists bool
	err := a.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM properties WHERE id = ? AND deleted_at IS NULL)", id).Scan(&exists)
	if err != nil || !exists {
		return nil, domainJob.ErrInvalidProperty
	}
	return true, nil
}

// GetWorkflowID obtiene el workflow_id del customer asociado a la propiedad
func (a *PropertyCheckerAdapter) GetWorkflowID(ctx context.Context, propertyID int64) (int64, error) {
	var workflowID int64
	err := a.db.QueryRowContext(ctx,
		`SELECT c.workflow_id FROM properties p
		 JOIN customers c ON c.id = p.customer_id
		 WHERE p.id = ? AND p.deleted_at IS NULL AND c.deleted_at IS NULL`,
		propertyID,
	).Scan(&workflowID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get workflow for property",
			slog.Int64("propertyId", propertyID),
			slog.String("error", err.Error()))
		return 0, domainJob.ErrInvalidWorkflow
	}
	slog.InfoContext(ctx, "Resolved workflow for property",
		slog.Int64("propertyId", propertyID),
		slog.Int64("workflowId", workflowID))
	return workflowID, nil
}

// UserCheckerAdapter adapta el repositorio de user para el checker del job use case
type UserCheckerAdapter struct {
	db *sql.DB
}

func NewUserCheckerAdapter(db *sql.DB) domainJob.UserChecker {
	return &UserCheckerAdapter{db: db}
}

func (a *UserCheckerAdapter) GetByID(ctx context.Context, id int64) (interface{}, error) {
	var exists bool
	err := a.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = ? AND is_active = 1 AND deleted_at IS NULL)", id).Scan(&exists)
	if err != nil || !exists {
		slog.ErrorContext(ctx, "User not found or inactive",
			slog.Int64("userId", id))
		return nil, domainJob.ErrInvalidUser
	}
	return true, nil
}

// TechnicianJobStatusCheckerAdapter adapta el repositorio de technician_job_status para el checker del job use case
type TechnicianJobStatusCheckerAdapter struct {
	db *sql.DB
}

func NewTechnicianJobStatusCheckerAdapter(db *sql.DB) domainJob.TechnicianJobStatusChecker {
	return &TechnicianJobStatusCheckerAdapter{db: db}
}

func (a *TechnicianJobStatusCheckerAdapter) GetByID(ctx context.Context, id int64) (interface{}, error) {
	var exists bool
	err := a.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM technician_job_statuses WHERE id = ? AND is_active = 1)", id).Scan(&exists)
	if err != nil || !exists {
		return nil, domainJob.ErrInvalidTechnicianJobStatus
	}
	return true, nil
}

// GetLinkedJobStatusID obtiene el job_status_id vinculado al technician_job_status
func (a *TechnicianJobStatusCheckerAdapter) GetLinkedJobStatusID(ctx context.Context, id int64) (*int64, error) {
	var jobStatusID sql.NullInt64
	err := a.db.QueryRowContext(ctx,
		"SELECT job_status_id FROM technician_job_statuses WHERE id = ? AND is_active = 1",
		id,
	).Scan(&jobStatusID)
	if err != nil {
		return nil, err
	}
	if !jobStatusID.Valid {
		return nil, nil
	}
	val := jobStatusID.Int64
	return &val, nil
}
