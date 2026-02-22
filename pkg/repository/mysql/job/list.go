package job

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	domainJob "github.com/your-org/jvairv2/pkg/domain/job"
)

// List obtiene una lista paginada de jobs con filtros opcionales
func (r *Repository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*domainJob.Job, int, error) {
	var conditions []string
	var args []interface{}

	// Siempre excluir soft-deleted
	conditions = append(conditions, "j.deleted_at IS NULL")

	// Filtro por closed (default: open jobs)
	if closed, ok := filters["closed"]; ok {
		switch v := closed.(type) {
		case string:
			switch v {
			case "all":
				// No filtrar por closed
			case "1", "true":
				conditions = append(conditions, "j.closed = 1")
			default:
				conditions = append(conditions, "j.closed = 0")
			}
		case bool:
			if v {
				conditions = append(conditions, "j.closed = 1")
			} else {
				conditions = append(conditions, "j.closed = 0")
			}
		}
	} else {
		// Default: solo jobs abiertos
		conditions = append(conditions, "j.closed = 0")
	}

	// Filtro por categoría
	if jobCategoryID, ok := filters["job_category_id"].(int64); ok && jobCategoryID > 0 {
		conditions = append(conditions, "j.job_category_id = ?")
		args = append(args, jobCategoryID)
	}

	// Filtro por estado
	if jobStatusID, ok := filters["job_status_id"].(int64); ok && jobStatusID > 0 {
		conditions = append(conditions, "j.job_status_id = ?")
		args = append(args, jobStatusID)
	}

	// Filtro por prioridad
	if jobPriorityID, ok := filters["job_priority_id"].(int64); ok && jobPriorityID > 0 {
		conditions = append(conditions, "j.job_priority_id = ?")
		args = append(args, jobPriorityID)
	}

	// Filtro por usuario
	if userID, ok := filters["user_id"]; ok {
		switch v := userID.(type) {
		case string:
			if v == "unassigned" {
				conditions = append(conditions, "j.user_id IS NULL")
			} else {
				conditions = append(conditions, "j.user_id = ?")
				args = append(args, v)
			}
		case int64:
			if v > 0 {
				conditions = append(conditions, "j.user_id = ?")
				args = append(args, v)
			}
		}
	}

	// Filtro por propiedad
	if propertyID, ok := filters["property_id"].(int64); ok && propertyID > 0 {
		conditions = append(conditions, "j.property_id = ?")
		args = append(args, propertyID)
	}

	// Filtro por workflow
	if workflowID, ok := filters["workflow_id"].(int64); ok && workflowID > 0 {
		conditions = append(conditions, "j.workflow_id = ?")
		args = append(args, workflowID)
	}

	// Búsqueda en múltiples campos (fiel al original: work_order, property fields, customer name)
	if search, ok := filters["search"].(string); ok && search != "" {
		searchCondition := `(
			j.work_order LIKE ? OR
			p.property_code LIKE ? OR
			p.street LIKE ? OR
			p.city LIKE ? OR
			p.state LIKE ? OR
			p.zip LIKE ? OR
			c.name LIKE ?
		)`
		conditions = append(conditions, searchCondition)
		searchPattern := "%" + search + "%"
		for i := 0; i < 7; i++ {
			args = append(args, searchPattern)
		}
	}

	whereClause := strings.Join(conditions, " AND ")

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM jobs j
		LEFT JOIN properties p ON p.id = j.property_id
		LEFT JOIN customers c ON c.id = p.customer_id
		WHERE %s
	`, whereClause)

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		slog.ErrorContext(ctx, "Failed to count jobs",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	// Sorting
	orderClause := "j.created_at DESC"
	if sort, ok := filters["sort"].(string); ok && sort != "" {
		direction := "DESC"
		if dir, ok := filters["direction"].(string); ok && strings.ToUpper(dir) == "ASC" {
			direction = "ASC"
		}

		switch sort {
		case "work_order":
			orderClause = fmt.Sprintf("j.work_order %s", direction)
		case "date_received":
			orderClause = fmt.Sprintf("j.date_received %s", direction)
		case "created_at":
			orderClause = fmt.Sprintf("j.created_at %s", direction)
		case "due_date":
			orderClause = fmt.Sprintf("j.due_date %s", direction)
		case "dispatch_date":
			orderClause = fmt.Sprintf("j.dispatch_date %s", direction)
		case "completion_date":
			orderClause = fmt.Sprintf("j.completion_date %s", direction)
		case "week_number":
			orderClause = fmt.Sprintf("j.week_number %s", direction)
		case "status":
			// Ordenar por workflow status order (fiel al original)
			orderClause = fmt.Sprintf("jsw.`order` %s", direction)
		}
	}

	// Data query
	offset := (page - 1) * pageSize
	dataQuery := fmt.Sprintf(`
		SELECT
			j.id, j.work_order, j.date_received, j.job_category_id, j.job_priority_id, j.job_status_id,
			j.technician_job_status_id, j.workflow_id, j.property_id, j.user_id, j.supervisor_ids,
			j.dispatch_date, j.completion_date, j.week_number, j.route_number,
			j.scheduled_time_type, j.scheduled_time, j.internal_job_notes, j.quick_notes,
			j.job_report, j.installation_due_date, j.cage_required, j.warranty_claim,
			j.warranty_registration, j.job_sales_price, j.money_turned_in, j.closed,
			j.dispatch_notes, j.call_logs, j.due_date, j.call_attempted,
			j.created_at, j.updated_at, j.deleted_at
		FROM jobs j
		LEFT JOIN properties p ON p.id = j.property_id
		LEFT JOIN customers c ON c.id = p.customer_id
		LEFT JOIN job_status_workflow jsw ON jsw.workflow_id = j.workflow_id AND jsw.job_status_id = j.job_status_id
		WHERE %s
		ORDER BY %s
		LIMIT ? OFFSET ?
	`, whereClause, orderClause)

	queryArgs := append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, queryArgs...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list jobs",
			slog.String("error", err.Error()))
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var jobs []*domainJob.Job
	for rows.Next() {
		j := &domainJob.Job{}
		if err := rows.Scan(
			&j.ID, &j.WorkOrder, &j.DateReceived, &j.JobCategoryID, &j.JobPriorityID, &j.JobStatusID,
			&j.TechnicianJobStatusID, &j.WorkflowID, &j.PropertyID, &j.UserID, &j.SupervisorIDs,
			&j.DispatchDate, &j.CompletionDate, &j.WeekNumber, &j.RouteNumber,
			&j.ScheduledTimeType, &j.ScheduledTime, &j.InternalJobNotes, &j.QuickNotes,
			&j.JobReport, &j.InstallationDueDate, &j.CageRequired, &j.WarrantyClaim,
			&j.WarrantyRegistration, &j.JobSalesPrice, &j.MoneyTurnedIn, &j.Closed,
			&j.DispatchNotes, &j.CallLogs, &j.DueDate, &j.CallAttempted,
			&j.CreatedAt, &j.UpdatedAt, &j.DeletedAt,
		); err != nil {
			slog.ErrorContext(ctx, "Failed to scan job row",
				slog.String("error", err.Error()))
			return nil, 0, err
		}
		jobs = append(jobs, j)
	}

	if err = rows.Err(); err != nil {
		slog.ErrorContext(ctx, "Error iterating job rows",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	return jobs, total, nil
}
