package job

import (
	"context"
	"database/sql"
	"log/slog"

	domainJob "github.com/angumol/jvairv2/pkg/domain/job"
)

// GetByID obtiene un job por su ID
func (r *Repository) GetByID(ctx context.Context, id int64) (*domainJob.Job, error) {
	query := `
		SELECT
			j.id, j.work_order, j.date_received, j.job_category_id, j.job_priority_id, j.job_status_id,
			j.technician_job_status_id, j.workflow_id, j.property_id, j.user_id, j.supervisor_ids,
			j.dispatch_date, j.completion_date, j.week_number, j.route_number,
			j.scheduled_time_type, j.scheduled_time, j.internal_job_notes, j.quick_notes,
			j.job_report, j.installation_due_date, j.cage_required, j.warranty_claim,
			j.warranty_registration, j.job_sales_price, j.money_turned_in, j.closed,
			j.dispatch_notes, j.call_logs, j.due_date, j.call_attempted,
			j.created_at, j.updated_at, j.deleted_at,
			jc.label as category_label, jc.type as category_type,
			js.label as status_label, js.class as status_class,
			jp.label as priority_label, jp.order as priority_order, jp.class as priority_class
		FROM jobs j
		INNER JOIN job_categories jc ON jc.id = j.job_category_id
		INNER JOIN job_statuses js ON js.id = j.job_status_id
		INNER JOIN job_priorities jp ON jp.id = j.job_priority_id
		WHERE j.id = ?
	`

	j := &domainJob.Job{}
	var categoryLabel, categoryType string
	var statusLabel string
	var statusClass *string
	var priorityLabel string
	var priorityOrder int
	var priorityClass *string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&j.ID, &j.WorkOrder, &j.DateReceived, &j.JobCategoryID, &j.JobPriorityID, &j.JobStatusID,
		&j.TechnicianJobStatusID, &j.WorkflowID, &j.PropertyID, &j.UserID, &j.SupervisorIDs,
		&j.DispatchDate, &j.CompletionDate, &j.WeekNumber, &j.RouteNumber,
		&j.ScheduledTimeType, &j.ScheduledTime, &j.InternalJobNotes, &j.QuickNotes,
		&j.JobReport, &j.InstallationDueDate, &j.CageRequired, &j.WarrantyClaim,
		&j.WarrantyRegistration, &j.JobSalesPrice, &j.MoneyTurnedIn, &j.Closed,
		&j.DispatchNotes, &j.CallLogs, &j.DueDate, &j.CallAttempted,
		&j.CreatedAt, &j.UpdatedAt, &j.DeletedAt,
		&categoryLabel, &categoryType,
		&statusLabel, &statusClass,
		&priorityLabel, &priorityOrder, &priorityClass,
	)
	if err == nil {
		// Asignar objetos anidados
		j.Category = &domainJob.JobCategory{
			ID:    j.JobCategoryID,
			Label: categoryLabel,
			Type:  categoryType,
		}
		j.Status = &domainJob.JobStatus{
			ID:    j.JobStatusID,
			Label: statusLabel,
			Class: statusClass,
		}
		j.Priority = &domainJob.JobPriority{
			ID:    j.JobPriorityID,
			Label: priorityLabel,
			Order: priorityOrder,
			Class: priorityClass,
		}
	}
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
