package new_dashboard

import (
	"context"
	"database/sql"
	"log/slog"

	domainNewDashboard "github.com/your-org/jvairv2/pkg/domain/new_dashboard"
)

// Repository implementa el repositorio MySQL para el dashboard enriquecido
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia del repositorio
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// buildTimeRangeCondition construye la condición SQL para el rango de tiempo
func buildTimeRangeCondition(timeRange domainNewDashboard.TimeRange, dateColumn string) (string, []interface{}) {
	switch timeRange {
	case domainNewDashboard.TimeRange7Days:
		return dateColumn + " >= DATE_SUB(NOW(), INTERVAL 7 DAY)", nil
	case domainNewDashboard.TimeRange30Days:
		return dateColumn + " >= DATE_SUB(NOW(), INTERVAL 30 DAY)", nil
	case domainNewDashboard.TimeRange90Days:
		return dateColumn + " >= DATE_SUB(NOW(), INTERVAL 90 DAY)", nil
	case domainNewDashboard.TimeRangeThisMonth:
		return dateColumn + " >= DATE_FORMAT(NOW(), '%Y-%m-01')", nil
	case domainNewDashboard.TimeRangeLastMonth:
		return dateColumn + " >= DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 1 MONTH), '%Y-%m-01') AND " +
			dateColumn + " < DATE_FORMAT(NOW(), '%Y-%m-01')", nil
	case domainNewDashboard.TimeRangeThisYear:
		return dateColumn + " >= DATE_FORMAT(NOW(), '%Y-01-01')", nil
	default:
		return dateColumn + " >= DATE_SUB(NOW(), INTERVAL 30 DAY)", nil
	}
}

// GetEnhancedStats obtiene estadísticas expandidas
func (r *Repository) GetEnhancedStats(ctx context.Context, userID *int64, timeRange domainNewDashboard.TimeRange) (*domainNewDashboard.EnhancedStats, error) {
	stats := &domainNewDashboard.EnhancedStats{}

	userFilter := ""
	var args []interface{}
	if userID != nil {
		userFilter = " AND j.user_id = ?"
		args = append(args, *userID)
	}

	// Jobs awaiting dispatch (todos los abiertos, sin filtro de tiempo)
	query := `SELECT COUNT(*) FROM jobs j WHERE j.deleted_at IS NULL AND j.closed = 0 AND j.job_status_id = 1` + userFilter
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&stats.JobsAwaitingDispatch); err != nil {
		slog.ErrorContext(ctx, "Failed to count jobs awaiting dispatch", slog.String("error", err.Error()))
		return nil, err
	}

	// Jobs dispatched (todos los abiertos, sin filtro de tiempo)
	query = `SELECT COUNT(*) FROM jobs j WHERE j.deleted_at IS NULL AND j.closed = 0 AND j.job_status_id = 2` + userFilter
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&stats.JobsDispatched); err != nil {
		slog.ErrorContext(ctx, "Failed to count dispatched jobs", slog.String("error", err.Error()))
		return nil, err
	}

	// Urgent jobs (todos los abiertos, sin filtro de tiempo)
	query = `SELECT COUNT(*) FROM jobs j WHERE j.deleted_at IS NULL AND j.closed = 0 AND j.job_priority_id = 4` + userFilter
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&stats.JobsUrgent); err != nil {
		slog.ErrorContext(ctx, "Failed to count urgent jobs", slog.String("error", err.Error()))
		return nil, err
	}

	// Total open jobs (todos los abiertos, sin filtro de tiempo)
	query = `SELECT COUNT(*) FROM jobs j WHERE j.deleted_at IS NULL AND j.closed = 0` + userFilter
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&stats.JobsOpen); err != nil {
		slog.ErrorContext(ctx, "Failed to count open jobs", slog.String("error", err.Error()))
		return nil, err
	}

	// Jobs closed (aplicar filtro de tiempo en updated_at)
	timeConditionClosed, _ := buildTimeRangeCondition(timeRange, "j.updated_at")
	query = `SELECT COUNT(*) FROM jobs j WHERE j.deleted_at IS NULL AND j.closed = 1 AND ` + timeConditionClosed + userFilter
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&stats.JobsClosedThisMonth); err != nil {
		slog.ErrorContext(ctx, "Failed to count jobs closed in period", slog.String("error", err.Error()))
		return nil, err
	}

	// Total alerts (filtrado por user_id si es técnico)
	alertUserFilter := ""
	if userID != nil {
		alertUserFilter = " AND a.user_id = ?"
	}
	query = `SELECT COUNT(*) FROM alerts a WHERE a.is_read = 0` + alertUserFilter
	if userID != nil {
		if err := r.db.QueryRowContext(ctx, query, *userID).Scan(&stats.TotalAlerts); err != nil {
			slog.ErrorContext(ctx, "Failed to count alerts", slog.String("error", err.Error()))
			return nil, err
		}
	} else {
		if err := r.db.QueryRowContext(ctx, query).Scan(&stats.TotalAlerts); err != nil {
			slog.ErrorContext(ctx, "Failed to count alerts", slog.String("error", err.Error()))
			return nil, err
		}
	}

	// Total tasks pending (filtrado por user_id si es técnico)
	taskUserFilter := ""
	if userID != nil {
		taskUserFilter = " AND jt.user_id = ?"
	}
	query = `SELECT COUNT(*) FROM job_tasks jt WHERE jt.deleted_at IS NULL AND jt.task_status_id IN (1, 2)` + taskUserFilter
	if userID != nil {
		if err := r.db.QueryRowContext(ctx, query, *userID).Scan(&stats.TotalTasksPending); err != nil {
			slog.ErrorContext(ctx, "Failed to count pending tasks", slog.String("error", err.Error()))
			return nil, err
		}
	} else {
		if err := r.db.QueryRowContext(ctx, query).Scan(&stats.TotalTasksPending); err != nil {
			slog.ErrorContext(ctx, "Failed to count pending tasks", slog.String("error", err.Error()))
			return nil, err
		}
	}

	// Total tasks overdue
	query = `SELECT COUNT(*) FROM job_tasks jt WHERE jt.deleted_at IS NULL AND jt.task_status_id IN (1, 2) AND jt.due_date < NOW()` + taskUserFilter
	if userID != nil {
		if err := r.db.QueryRowContext(ctx, query, *userID).Scan(&stats.TotalTasksOverdue); err != nil {
			slog.ErrorContext(ctx, "Failed to count overdue tasks", slog.String("error", err.Error()))
			return nil, err
		}
	} else {
		if err := r.db.QueryRowContext(ctx, query).Scan(&stats.TotalTasksOverdue); err != nil {
			slog.ErrorContext(ctx, "Failed to count overdue tasks", slog.String("error", err.Error()))
			return nil, err
		}
	}

	// Solo para admin: conteos de invoices, quotes, warranties
	if userID == nil {
		// Total invoices pending (invoices no tiene status, contar todas las no eliminadas)
		query = `SELECT COUNT(*) FROM invoices i WHERE i.deleted_at IS NULL`
		if err := r.db.QueryRowContext(ctx, query).Scan(&stats.TotalInvoicesPending); err != nil {
			slog.ErrorContext(ctx, "Failed to count pending invoices", slog.String("error", err.Error()))
			return nil, err
		}

		// Total quotes pending
		query = `SELECT COUNT(*) FROM quotes q WHERE q.deleted_at IS NULL AND q.quote_status_id = 1`
		if err := r.db.QueryRowContext(ctx, query).Scan(&stats.TotalQuotesPending); err != nil {
			slog.ErrorContext(ctx, "Failed to count pending quotes", slog.String("error", err.Error()))
			return nil, err
		}

		// Total warranty claims
		query = `SELECT COUNT(*) FROM warranty_claims wc WHERE wc.deleted_at IS NULL AND wc.warranty_claim_status_id IN (1, 2)`
		if err := r.db.QueryRowContext(ctx, query).Scan(&stats.TotalWarrantyClaims); err != nil {
			slog.ErrorContext(ctx, "Failed to count warranty claims", slog.String("error", err.Error()))
			return nil, err
		}
	}

	return stats, nil
}

// GetAlertSummary obtiene el resumen de alertas
func (r *Repository) GetAlertSummary(ctx context.Context, userID *int64, timeRange domainNewDashboard.TimeRange) (*domainNewDashboard.AlertSummary, error) {
	summary := &domainNewDashboard.AlertSummary{}

	userFilter := ""
	var args []interface{}
	if userID != nil {
		userFilter = " AND a.user_id = ?"
		args = append(args, *userID)
	}

	// Count unread alerts
	query := `SELECT COUNT(*) FROM alerts a WHERE a.is_read = 0` + userFilter
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&summary.UnreadCount); err != nil {
		slog.ErrorContext(ctx, "Failed to count unread alerts", slog.String("error", err.Error()))
		return nil, err
	}

	// Get recent alerts (limit 10)
	query = `SELECT a.id, a.alert_type, a.message, a.message_level, a.created_at
		FROM alerts a
		WHERE a.is_read = 0` + userFilter + `
		ORDER BY a.created_at DESC
		LIMIT 10`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get alerts", slog.String("error", err.Error()))
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.ErrorContext(ctx, "Failed to close rows", slog.String("error", err.Error()))
		}
	}()

	var alerts []*domainNewDashboard.Alert
	for rows.Next() {
		var alert domainNewDashboard.Alert
		if err := rows.Scan(&alert.ID, &alert.AlertType, &alert.Message, &alert.MessageLevel, &alert.CreatedAt); err != nil {
			slog.ErrorContext(ctx, "Failed to scan alert", slog.String("error", err.Error()))
			return nil, err
		}
		alerts = append(alerts, &alert)
	}

	if alerts == nil {
		alerts = []*domainNewDashboard.Alert{}
	}
	summary.Alerts = alerts

	return summary, nil
}

// GetTaskSummary obtiene el resumen de tareas
func (r *Repository) GetTaskSummary(ctx context.Context, userID *int64, timeRange domainNewDashboard.TimeRange) (*domainNewDashboard.TaskSummary, error) {
	summary := &domainNewDashboard.TaskSummary{}

	userFilter := ""
	var args []interface{}
	if userID != nil {
		userFilter = " AND jt.user_id = ?"
		args = append(args, *userID)
	}

	// Count pending tasks
	query := `SELECT COUNT(*) FROM job_tasks jt WHERE jt.deleted_at IS NULL AND jt.task_status_id IN (1, 2)` + userFilter
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&summary.TotalPending); err != nil {
		slog.ErrorContext(ctx, "Failed to count pending tasks", slog.String("error", err.Error()))
		return nil, err
	}

	// Count overdue tasks
	query = `SELECT COUNT(*) FROM job_tasks jt WHERE jt.deleted_at IS NULL AND jt.task_status_id IN (1, 2) AND jt.due_date < NOW()` + userFilter
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&summary.TotalOverdue); err != nil {
		slog.ErrorContext(ctx, "Failed to count overdue tasks", slog.String("error", err.Error()))
		return nil, err
	}

	// Get recent tasks (limit 10)
	query = `SELECT jt.id, jt.job_id, j.work_order, jt.task, jt.due_date, ts.label,
			CASE WHEN jt.due_date < NOW() THEN 1 ELSE 0 END as is_overdue
		FROM job_tasks jt
		JOIN jobs j ON j.id = jt.job_id
		JOIN task_statuses ts ON ts.id = jt.task_status_id
		WHERE jt.deleted_at IS NULL AND jt.task_status_id IN (1, 2)` + userFilter + `
		ORDER BY jt.due_date ASC
		LIMIT 10`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get tasks", slog.String("error", err.Error()))
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.ErrorContext(ctx, "Failed to close rows", slog.String("error", err.Error()))
		}
	}()

	var tasks []*domainNewDashboard.Task
	for rows.Next() {
		var task domainNewDashboard.Task
		var workOrder sql.NullString
		if err := rows.Scan(&task.ID, &task.JobID, &workOrder, &task.Title, &task.DueDate, &task.StatusName, &task.IsOverdue); err != nil {
			slog.ErrorContext(ctx, "Failed to scan task", slog.String("error", err.Error()))
			return nil, err
		}
		if workOrder.Valid {
			task.WorkOrder = &workOrder.String
		}
		tasks = append(tasks, &task)
	}

	if tasks == nil {
		tasks = []*domainNewDashboard.Task{}
	}
	summary.Tasks = tasks

	return summary, nil
}

// GetRecentActivity obtiene actividad reciente
func (r *Repository) GetRecentActivity(ctx context.Context, userID *int64, timeRange domainNewDashboard.TimeRange) ([]*domainNewDashboard.Activity, error) {
	userFilter := ""
	var args []interface{}
	if userID != nil {
		userFilter = " AND j.user_id = ?"
		args = append(args, *userID)
	}

	query := `SELECT jal.id, jal.job_id, j.work_order, jal.type, jal.log, u.name, jal.created_at
		FROM job_activity_logs jal
		JOIN jobs j ON j.id = jal.job_id
		LEFT JOIN users u ON u.id = jal.user_id
		WHERE jal.created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)` + userFilter + `
		ORDER BY jal.created_at DESC
		LIMIT 15`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get recent activity", slog.String("error", err.Error()))
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.ErrorContext(ctx, "Failed to close rows", slog.String("error", err.Error()))
		}
	}()

	var activities []*domainNewDashboard.Activity
	for rows.Next() {
		var activity domainNewDashboard.Activity
		var workOrder, userName sql.NullString
		if err := rows.Scan(&activity.ID, &activity.JobID, &workOrder, &activity.Type, &activity.Log, &userName, &activity.CreatedAt); err != nil {
			slog.ErrorContext(ctx, "Failed to scan activity", slog.String("error", err.Error()))
			return nil, err
		}
		if workOrder.Valid {
			activity.WorkOrder = &workOrder.String
		}
		activity.UserName = userName.String
		activities = append(activities, &activity)
	}

	if activities == nil {
		activities = []*domainNewDashboard.Activity{}
	}

	return activities, nil
}

// GetInvoiceSummary obtiene el resumen de facturación (solo admin)
func (r *Repository) GetInvoiceSummary(ctx context.Context, timeRange domainNewDashboard.TimeRange) (*domainNewDashboard.InvoiceSummary, error) {
	summary := &domainNewDashboard.InvoiceSummary{}

	timeCondition, _ := buildTimeRangeCondition(timeRange, "i.created_at")

	// Count pending invoices (invoices no tiene status, contar todas)
	query := `SELECT COUNT(*), COALESCE(SUM(i.total), 0) FROM invoices i
		WHERE i.deleted_at IS NULL AND ` + timeCondition
	if err := r.db.QueryRowContext(ctx, query).Scan(&summary.TotalPending, &summary.AmountPending); err != nil {
		slog.ErrorContext(ctx, "Failed to count pending invoices", slog.String("error", err.Error()))
		return nil, err
	}

	// Count paid invoices (usar invoice_payments para determinar pagadas)
	query = `SELECT COUNT(DISTINCT i.id) FROM invoices i
		INNER JOIN invoice_payments ip ON ip.invoice_id = i.id AND ip.status = 'paid'
		WHERE i.deleted_at IS NULL AND ` + timeCondition
	if err := r.db.QueryRowContext(ctx, query).Scan(&summary.TotalPaid); err != nil {
		slog.ErrorContext(ctx, "Failed to count paid invoices", slog.String("error", err.Error()))
		return nil, err
	}

	// Count overdue invoices (invoices no tiene due_date, omitir este conteo)
	summary.TotalOverdue = 0

	// Amount paid this month
	query = `SELECT COALESCE(SUM(ip.amount), 0) FROM invoice_payments ip
		WHERE ip.created_at >= DATE_FORMAT(NOW(), '%Y-%m-01') AND ip.status = 'paid'`
	if err := r.db.QueryRowContext(ctx, query).Scan(&summary.AmountPaidThisMonth); err != nil {
		slog.ErrorContext(ctx, "Failed to sum paid amount", slog.String("error", err.Error()))
		return nil, err
	}

	// Get recent invoices (limit 5)
	query = `SELECT i.id, i.invoice_number, i.job_id, c.name, i.total,
			CASE WHEN ip.id IS NOT NULL THEN 'paid' ELSE 'pending' END as status,
			i.created_at
		FROM invoices i
		LEFT JOIN jobs j ON j.id = i.job_id
		LEFT JOIN properties p ON p.id = j.property_id
		LEFT JOIN customers c ON c.id = p.customer_id
		LEFT JOIN invoice_payments ip ON ip.invoice_id = i.id AND ip.status = 'paid'
		WHERE i.deleted_at IS NULL AND ` + timeCondition + `
		ORDER BY i.created_at DESC
		LIMIT 5`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get recent invoices", slog.String("error", err.Error()))
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.ErrorContext(ctx, "Failed to close rows", slog.String("error", err.Error()))
		}
	}()

	var invoices []*domainNewDashboard.RecentInvoice
	for rows.Next() {
		var invoice domainNewDashboard.RecentInvoice
		var jobID sql.NullInt64
		var customerName sql.NullString
		if err := rows.Scan(&invoice.ID, &invoice.InvoiceNumber, &jobID, &customerName, &invoice.Total, &invoice.Status, &invoice.CreatedAt); err != nil {
			slog.ErrorContext(ctx, "Failed to scan invoice", slog.String("error", err.Error()))
			return nil, err
		}
		if jobID.Valid {
			invoice.JobID = &jobID.Int64
		}
		invoice.CustomerName = customerName.String
		invoices = append(invoices, &invoice)
	}

	if invoices == nil {
		invoices = []*domainNewDashboard.RecentInvoice{}
	}
	summary.RecentInvoices = invoices

	return summary, nil
}

// GetQuoteSummary obtiene el resumen de cotizaciones (solo admin)
func (r *Repository) GetQuoteSummary(ctx context.Context, timeRange domainNewDashboard.TimeRange) (*domainNewDashboard.QuoteSummary, error) {
	summary := &domainNewDashboard.QuoteSummary{}

	timeCondition, _ := buildTimeRangeCondition(timeRange, "q.created_at")

	// Count pending quotes
	query := `SELECT COUNT(*) FROM quotes q WHERE q.deleted_at IS NULL AND q.quote_status_id = 1 AND ` + timeCondition
	if err := r.db.QueryRowContext(ctx, query).Scan(&summary.TotalPending); err != nil {
		slog.ErrorContext(ctx, "Failed to count pending quotes", slog.String("error", err.Error()))
		return nil, err
	}

	// Count approved quotes
	query = `SELECT COUNT(*) FROM quotes q WHERE q.deleted_at IS NULL AND q.quote_status_id = 2 AND ` + timeCondition
	if err := r.db.QueryRowContext(ctx, query).Scan(&summary.TotalApproved); err != nil {
		slog.ErrorContext(ctx, "Failed to count approved quotes", slog.String("error", err.Error()))
		return nil, err
	}

	// Count rejected quotes
	query = `SELECT COUNT(*) FROM quotes q WHERE q.deleted_at IS NULL AND q.quote_status_id = 3 AND ` + timeCondition
	if err := r.db.QueryRowContext(ctx, query).Scan(&summary.TotalRejected); err != nil {
		slog.ErrorContext(ctx, "Failed to count rejected quotes", slog.String("error", err.Error()))
		return nil, err
	}

	// Get recent quotes (limit 5)
	query = `SELECT q.id, q.quote_number, q.job_id, c.name, q.amount, qs.label, q.created_at
		FROM quotes q
		LEFT JOIN jobs j ON j.id = q.job_id
		LEFT JOIN properties p ON p.id = j.property_id
		LEFT JOIN customers c ON c.id = p.customer_id
		JOIN quote_statuses qs ON qs.id = q.quote_status_id
		WHERE q.deleted_at IS NULL AND ` + timeCondition + `
		ORDER BY q.created_at DESC
		LIMIT 5`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get recent quotes", slog.String("error", err.Error()))
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.ErrorContext(ctx, "Failed to close rows", slog.String("error", err.Error()))
		}
	}()

	var quotes []*domainNewDashboard.RecentQuote
	for rows.Next() {
		var quote domainNewDashboard.RecentQuote
		var jobID sql.NullInt64
		var customerName sql.NullString
		if err := rows.Scan(&quote.ID, &quote.QuoteNumber, &jobID, &customerName, &quote.Total, &quote.StatusName, &quote.CreatedAt); err != nil {
			slog.ErrorContext(ctx, "Failed to scan quote", slog.String("error", err.Error()))
			return nil, err
		}
		if jobID.Valid {
			quote.JobID = &jobID.Int64
		}
		quote.CustomerName = customerName.String
		quotes = append(quotes, &quote)
	}

	if quotes == nil {
		quotes = []*domainNewDashboard.RecentQuote{}
	}
	summary.RecentQuotes = quotes

	return summary, nil
}

// GetWarrantySummary obtiene el resumen de garantías (solo admin)
func (r *Repository) GetWarrantySummary(ctx context.Context, timeRange domainNewDashboard.TimeRange) (*domainNewDashboard.WarrantySummary, error) {
	summary := &domainNewDashboard.WarrantySummary{}

	// Count active warranties (warranties no tiene end_date, contar por status activo)
	query := `SELECT COUNT(*) FROM warranties w WHERE w.deleted_at IS NULL AND w.warranty_status_id IN (1, 2)`
	if err := r.db.QueryRowContext(ctx, query).Scan(&summary.ActiveWarranties); err != nil {
		slog.ErrorContext(ctx, "Failed to count active warranties", slog.String("error", err.Error()))
		return nil, err
	}

	// Count warranties expiring this month (warranties no tiene end_date, omitir)
	summary.ExpiringThisMonth = 0

	// Count open warranty claims
	query = `SELECT COUNT(*) FROM warranty_claims wc WHERE wc.deleted_at IS NULL AND wc.warranty_claim_status_id IN (1, 2)`
	if err := r.db.QueryRowContext(ctx, query).Scan(&summary.OpenClaims); err != nil {
		slog.ErrorContext(ctx, "Failed to count warranty claims", slog.String("error", err.Error()))
		return nil, err
	}

	// Get recent claims (limit 5)
	query = `SELECT wc.id, wc.claim_number, wc.job_id, c.name, wcs.label, wc.created_at
		FROM warranty_claims wc
		LEFT JOIN jobs j ON j.id = wc.job_id
		LEFT JOIN properties p ON p.id = j.property_id
		LEFT JOIN customers c ON c.id = p.customer_id
		JOIN warranty_claim_statuses wcs ON wcs.id = wc.warranty_claim_status_id
		WHERE wc.deleted_at IS NULL
		ORDER BY wc.created_at DESC
		LIMIT 5`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get recent warranty claims", slog.String("error", err.Error()))
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.ErrorContext(ctx, "Failed to close rows", slog.String("error", err.Error()))
		}
	}()

	var claims []*domainNewDashboard.RecentClaim
	for rows.Next() {
		var claim domainNewDashboard.RecentClaim
		var jobID sql.NullInt64
		var customerName sql.NullString
		if err := rows.Scan(&claim.ID, &claim.ClaimNumber, &jobID, &customerName, &claim.StatusName, &claim.CreatedAt); err != nil {
			slog.ErrorContext(ctx, "Failed to scan warranty claim", slog.String("error", err.Error()))
			return nil, err
		}
		if jobID.Valid {
			claim.JobID = &jobID.Int64
		}
		claim.CustomerName = customerName.String
		claims = append(claims, &claim)
	}

	if claims == nil {
		claims = []*domainNewDashboard.RecentClaim{}
	}
	summary.RecentClaims = claims

	return summary, nil
}

// GetJobsByCategory obtiene la distribución de jobs por categoría (solo admin)
func (r *Repository) GetJobsByCategory(ctx context.Context, timeRange domainNewDashboard.TimeRange) ([]*domainNewDashboard.CategoryCount, error) {
	// Mostrar distribución de todos los jobs abiertos (sin filtro de tiempo)
	query := `SELECT jc.id, jc.label, COUNT(j.id) as count
		FROM jobs j
		JOIN job_categories jc ON jc.id = j.job_category_id
		WHERE j.deleted_at IS NULL AND j.closed = 0
		GROUP BY jc.id, jc.label
		ORDER BY count DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get jobs by category", slog.String("error", err.Error()))
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.ErrorContext(ctx, "Failed to close rows", slog.String("error", err.Error()))
		}
	}()

	var categories []*domainNewDashboard.CategoryCount
	for rows.Next() {
		var cat domainNewDashboard.CategoryCount
		if err := rows.Scan(&cat.CategoryID, &cat.CategoryName, &cat.Count); err != nil {
			slog.ErrorContext(ctx, "Failed to scan category count", slog.String("error", err.Error()))
			return nil, err
		}
		categories = append(categories, &cat)
	}

	if categories == nil {
		categories = []*domainNewDashboard.CategoryCount{}
	}

	return categories, nil
}

// GetJobsByStatus obtiene la distribución de jobs por estado
func (r *Repository) GetJobsByStatus(ctx context.Context, userID *int64, timeRange domainNewDashboard.TimeRange) ([]*domainNewDashboard.StatusCount, error) {
	// Mostrar distribución de todos los jobs (abiertos y cerrados sin filtro de tiempo)
	userFilter := ""
	var args []interface{}
	if userID != nil {
		userFilter = " AND j.user_id = ?"
		args = append(args, *userID)
	}

	query := `SELECT js.id, js.label, COUNT(j.id) as count
		FROM jobs j
		JOIN job_statuses js ON js.id = j.job_status_id
		WHERE j.deleted_at IS NULL` + userFilter + `
		GROUP BY js.id, js.label
		ORDER BY count DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get jobs by status", slog.String("error", err.Error()))
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.ErrorContext(ctx, "Failed to close rows", slog.String("error", err.Error()))
		}
	}()

	var statuses []*domainNewDashboard.StatusCount
	for rows.Next() {
		var status domainNewDashboard.StatusCount
		if err := rows.Scan(&status.StatusID, &status.StatusName, &status.Count); err != nil {
			slog.ErrorContext(ctx, "Failed to scan status count", slog.String("error", err.Error()))
			return nil, err
		}
		statuses = append(statuses, &status)
	}

	if statuses == nil {
		statuses = []*domainNewDashboard.StatusCount{}
	}

	return statuses, nil
}

// GetJobsDueThisWeek obtiene jobs que vencen esta semana
func (r *Repository) GetJobsDueThisWeek(ctx context.Context, userID *int64) ([]*domainNewDashboard.DueJob, error) {
	userFilter := ""
	var args []interface{}
	if userID != nil {
		userFilter = " AND j.user_id = ?"
		args = append(args, *userID)
	}

	query := `SELECT j.id, j.work_order, j.due_date, c.name, p.street, js.label, jp.label
		FROM jobs j
		JOIN properties p ON p.id = j.property_id
		JOIN customers c ON c.id = p.customer_id
		JOIN job_statuses js ON js.id = j.job_status_id
		JOIN job_priorities jp ON jp.id = j.job_priority_id
		WHERE j.deleted_at IS NULL
		AND j.closed = 0
		AND j.due_date >= CURDATE()
		AND j.due_date < DATE_ADD(CURDATE(), INTERVAL 7 DAY)` + userFilter + `
		ORDER BY j.due_date ASC
		LIMIT 10`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get jobs due this week", slog.String("error", err.Error()))
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.ErrorContext(ctx, "Failed to close rows", slog.String("error", err.Error()))
		}
	}()

	var jobs []*domainNewDashboard.DueJob
	for rows.Next() {
		var job domainNewDashboard.DueJob
		var workOrder sql.NullString
		if err := rows.Scan(&job.ID, &workOrder, &job.DueDate, &job.CustomerName, &job.PropertyStreet, &job.StatusName, &job.PriorityName); err != nil {
			slog.ErrorContext(ctx, "Failed to scan due job", slog.String("error", err.Error()))
			return nil, err
		}
		if workOrder.Valid {
			job.WorkOrder = &workOrder.String
		}
		jobs = append(jobs, &job)
	}

	if jobs == nil {
		jobs = []*domainNewDashboard.DueJob{}
	}

	return jobs, nil
}

// GetTechnicianWorkload obtiene la carga de trabajo por técnico (solo admin)
func (r *Repository) GetTechnicianWorkload(ctx context.Context, timeRange domainNewDashboard.TimeRange) ([]*domainNewDashboard.TechnicianLoad, error) {
	timeCondition, _ := buildTimeRangeCondition(timeRange, "j.date_received")

	query := `SELECT u.id, u.name,
		COUNT(CASE WHEN j.closed = 0 THEN 1 END) as open_jobs,
		COUNT(CASE WHEN j.closed = 0 AND j.job_priority_id = 4 THEN 1 END) as urgent_jobs
		FROM users u
		JOIN assigned_roles ar ON ar.entity_id = u.id AND ar.entity_type = 'App\\Models\\User'
		JOIN roles r ON r.id = ar.role_id AND r.name IN ('technician', 'installer')
		LEFT JOIN jobs j ON j.user_id = u.id AND j.deleted_at IS NULL AND ` + timeCondition + `
		WHERE u.deleted_at IS NULL AND u.is_active = 1
		GROUP BY u.id, u.name
		ORDER BY open_jobs DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get technician workload", slog.String("error", err.Error()))
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.ErrorContext(ctx, "Failed to close rows", slog.String("error", err.Error()))
		}
	}()

	var workloads []*domainNewDashboard.TechnicianLoad
	for rows.Next() {
		var load domainNewDashboard.TechnicianLoad
		if err := rows.Scan(&load.UserID, &load.Name, &load.OpenJobs, &load.UrgentJobs); err != nil {
			slog.ErrorContext(ctx, "Failed to scan technician load", slog.String("error", err.Error()))
			return nil, err
		}
		workloads = append(workloads, &load)
	}

	if workloads == nil {
		workloads = []*domainNewDashboard.TechnicianLoad{}
	}

	return workloads, nil
}
