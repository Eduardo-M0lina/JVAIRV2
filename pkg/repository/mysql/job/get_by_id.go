package job

import (
	"context"
	"database/sql"
	"log/slog"

	domainJob "github.com/your-org/jvairv2/pkg/domain/job"
)

// GetByID obtiene un job por su ID
func (r *Repository) GetByID(ctx context.Context, id int64) (*domainJob.Job, error) {
	query := `
		SELECT
			id, work_order, date_received, job_category_id, job_priority_id, job_status_id,
			technician_job_status_id, workflow_id, property_id, user_id, supervisor_ids,
			dispatch_date, completion_date, week_number, route_number,
			scheduled_time_type, scheduled_time, internal_job_notes, quick_notes,
			job_report, installation_due_date, cage_required, warranty_claim,
			warranty_registration, job_sales_price, money_turned_in, closed,
			dispatch_notes, call_logs, due_date, call_attempted,
			created_at, updated_at, deleted_at
		FROM jobs
		WHERE id = ?
	`

	j := &domainJob.Job{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&j.ID, &j.WorkOrder, &j.DateReceived, &j.JobCategoryID, &j.JobPriorityID, &j.JobStatusID,
		&j.TechnicianJobStatusID, &j.WorkflowID, &j.PropertyID, &j.UserID, &j.SupervisorIDs,
		&j.DispatchDate, &j.CompletionDate, &j.WeekNumber, &j.RouteNumber,
		&j.ScheduledTimeType, &j.ScheduledTime, &j.InternalJobNotes, &j.QuickNotes,
		&j.JobReport, &j.InstallationDueDate, &j.CageRequired, &j.WarrantyClaim,
		&j.WarrantyRegistration, &j.JobSalesPrice, &j.MoneyTurnedIn, &j.Closed,
		&j.DispatchNotes, &j.CallLogs, &j.DueDate, &j.CallAttempted,
		&j.CreatedAt, &j.UpdatedAt, &j.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domainJob.ErrJobNotFound
		}
		slog.ErrorContext(ctx, "Failed to get job by ID",
			slog.Int64("id", id),
			slog.String("error", err.Error()))
		return nil, err
	}

	return j, nil
}
