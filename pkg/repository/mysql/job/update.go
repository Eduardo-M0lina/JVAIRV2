package job

import (
	"context"
	"log/slog"

	domainJob "github.com/your-org/jvairv2/pkg/domain/job"
)

// Update actualiza un job existente
func (r *Repository) Update(ctx context.Context, j *domainJob.Job) error {
	query := `
		UPDATE jobs SET
			work_order = ?, date_received = ?, job_category_id = ?, job_priority_id = ?, job_status_id = ?,
			technician_job_status_id = ?, workflow_id = ?, property_id = ?, user_id = ?, supervisor_ids = ?,
			dispatch_date = ?, completion_date = ?, week_number = ?, route_number = ?,
			scheduled_time_type = ?, scheduled_time = ?, internal_job_notes = ?, quick_notes = ?,
			job_report = ?, installation_due_date = ?, cage_required = ?, warranty_claim = ?,
			warranty_registration = ?, job_sales_price = ?, money_turned_in = ?, closed = ?,
			dispatch_notes = ?, call_logs = ?, due_date = ?, call_attempted = ?,
			updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query,
		j.WorkOrder, j.DateReceived, j.JobCategoryID, j.JobPriorityID, j.JobStatusID,
		j.TechnicianJobStatusID, j.WorkflowID, j.PropertyID, j.UserID, j.SupervisorIDs,
		j.DispatchDate, j.CompletionDate, j.WeekNumber, j.RouteNumber,
		j.ScheduledTimeType, j.ScheduledTime, j.InternalJobNotes, j.QuickNotes,
		j.JobReport, j.InstallationDueDate, j.CageRequired, j.WarrantyClaim,
		j.WarrantyRegistration, j.JobSalesPrice, j.MoneyTurnedIn, j.Closed,
		j.DispatchNotes, j.CallLogs, j.DueDate, j.CallAttempted,
		j.ID,
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to update job",
			slog.Int64("id", j.ID),
			slog.String("error", err.Error()))
		return err
	}

	return nil
}
