package dashboard

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	domainDashboard "github.com/your-org/jvairv2/pkg/domain/dashboard"
)

const timeFormat = "2006-01-02T15:04:05Z07:00"

// Repository implementa el repositorio MySQL para el dashboard
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia del repositorio de dashboard
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// GetStats obtiene las estadísticas generales del dashboard aplicando filtros
func (r *Repository) GetStats(ctx context.Context, userID *int64, filters domainDashboard.DashboardFilters) (*domainDashboard.DashboardStats, error) {
	stats := &domainDashboard.DashboardStats{}

	// Construir filtros base
	var whereClauses []string
	var args []interface{}

	whereClauses = append(whereClauses, "j.deleted_at IS NULL")

	if userID != nil {
		whereClauses = append(whereClauses, "j.user_id = ?")
		args = append(args, *userID)
	}

	// Aplicar filtros adicionales del usuario
	if filters.UserID != nil {
		if *filters.UserID == -1 {
			whereClauses = append(whereClauses, "j.user_id IS NULL")
		} else {
			whereClauses = append(whereClauses, "j.user_id = ?")
			args = append(args, *filters.UserID)
		}
	}

	if filters.JobStatusID != nil {
		whereClauses = append(whereClauses, "j.job_status_id = ?")
		args = append(args, *filters.JobStatusID)
	}

	if filters.JobPriorityID != nil {
		whereClauses = append(whereClauses, "j.job_priority_id = ?")
		args = append(args, *filters.JobPriorityID)
	}

	if filters.Year != nil {
		whereClauses = append(whereClauses, "YEAR(j.date_received) = ?")
		args = append(args, *filters.Year)
	}

	if filters.Week != nil {
		whereClauses = append(whereClauses, "j.week_number = ?")
		args = append(args, *filters.Week)
	}

	if filters.LastDays != nil && *filters.LastDays > 0 {
		whereClauses = append(whereClauses, "j.date_received >= DATE_SUB(NOW(), INTERVAL ? DAY)")
		args = append(args, *filters.LastDays)
	}

	// Búsqueda de texto (si aplica)
	searchArgs := []interface{}{}
	searchClause := ""
	if filters.Search != nil && *filters.Search != "" {
		searchPattern := "%" + *filters.Search + "%"
		searchClause = ` AND (
			j.work_order LIKE ? OR
			p.property_code LIKE ? OR
			p.street LIKE ? OR
			p.city LIKE ? OR
			p.state LIKE ? OR
			p.zip LIKE ? OR
			c.name LIKE ? OR
			c.email LIKE ? OR
			c.phone LIKE ? OR
			c.mobile LIKE ?
		)`
		for i := 0; i < 10; i++ {
			searchArgs = append(searchArgs, searchPattern)
		}
	}

	// Construir WHERE base
	whereSQL := " WHERE " + whereClauses[0]
	for i := 1; i < len(whereClauses); i++ {
		whereSQL += " AND " + whereClauses[i]
	}

	// Jobs awaiting dispatch (job_status_id = 1, not closed)
	query := `SELECT COUNT(*) FROM jobs j
		LEFT JOIN properties p ON p.id = j.property_id
		LEFT JOIN customers c ON c.id = p.customer_id` +
		whereSQL + ` AND j.closed = 0 AND j.job_status_id = 1` + searchClause
	countArgs := append(append([]interface{}{}, args...), searchArgs...)
	if err := r.db.QueryRowContext(ctx, query, countArgs...).Scan(&stats.JobsAwaitingDispatch); err != nil {
		slog.ErrorContext(ctx, "Failed to count jobs awaiting dispatch", slog.String("error", err.Error()))
		return nil, err
	}

	// Jobs dispatched (job_status_id = 2, not closed)
	query = `SELECT COUNT(*) FROM jobs j
		LEFT JOIN properties p ON p.id = j.property_id
		LEFT JOIN customers c ON c.id = p.customer_id` +
		whereSQL + ` AND j.closed = 0 AND j.job_status_id = 2` + searchClause
	countArgs = append(append([]interface{}{}, args...), searchArgs...)
	if err := r.db.QueryRowContext(ctx, query, countArgs...).Scan(&stats.JobsDispatched); err != nil {
		slog.ErrorContext(ctx, "Failed to count dispatched jobs", slog.String("error", err.Error()))
		return nil, err
	}

	// Urgent jobs (job_priority_id = 4, not closed)
	query = `SELECT COUNT(*) FROM jobs j
		LEFT JOIN properties p ON p.id = j.property_id
		LEFT JOIN customers c ON c.id = p.customer_id` +
		whereSQL + ` AND j.closed = 0 AND j.job_priority_id = 4` + searchClause
	countArgs = append(append([]interface{}{}, args...), searchArgs...)
	if err := r.db.QueryRowContext(ctx, query, countArgs...).Scan(&stats.JobsUrgent); err != nil {
		slog.ErrorContext(ctx, "Failed to count urgent jobs", slog.String("error", err.Error()))
		return nil, err
	}

	// Total open jobs
	openWhereSQL := whereSQL
	if filters.IsClosed == nil || *filters.IsClosed != "all" {
		openWhereSQL += " AND j.closed = 0"
	} else if filters.IsClosed != nil && *filters.IsClosed == "1" {
		openWhereSQL += " AND j.closed = 1"
	}
	query = `SELECT COUNT(*) FROM jobs j
		LEFT JOIN properties p ON p.id = j.property_id
		LEFT JOIN customers c ON c.id = p.customer_id` +
		openWhereSQL + searchClause
	countArgs = append(append([]interface{}{}, args...), searchArgs...)
	if err := r.db.QueryRowContext(ctx, query, countArgs...).Scan(&stats.JobsOpen); err != nil {
		slog.ErrorContext(ctx, "Failed to count open jobs", slog.String("error", err.Error()))
		return nil, err
	}

	// Jobs closed this month
	firstOfMonth := time.Now().UTC().Format("2006-01-01")
	query = `SELECT COUNT(*) FROM jobs j
		LEFT JOIN properties p ON p.id = j.property_id
		LEFT JOIN customers c ON c.id = p.customer_id` +
		whereSQL + ` AND j.closed = 1 AND j.updated_at >= ?` + searchClause
	monthArgs := append([]interface{}{firstOfMonth}, args...)
	monthArgs = append(monthArgs, searchArgs...)
	if err := r.db.QueryRowContext(ctx, query, monthArgs...).Scan(&stats.JobsClosedThisMonth); err != nil {
		slog.ErrorContext(ctx, "Failed to count jobs closed this month", slog.String("error", err.Error()))
		return nil, err
	}

	return stats, nil
}

// dashboardJobQuery retorna el SELECT base con JOINs enriquecidos para el dashboard
func dashboardJobQuery() string {
	return `
		SELECT
			j.id, j.work_order, j.date_received, j.job_category_id, j.job_priority_id,
			j.job_status_id, j.property_id, j.user_id, j.dispatch_date, j.due_date,
			j.quick_notes, j.closed, j.call_attempted, j.cage_required, j.week_number,
			j.route_number, j.job_sales_price, j.scheduled_time, j.completion_date,
			p.street, p.city, p.zip,
			c.name,
			u.name,
			jc.label,
			jp.label, jp.class,
			js.label, js.class
		FROM jobs j
		LEFT JOIN properties p ON p.id = j.property_id
		LEFT JOIN customers c ON c.id = p.customer_id
		LEFT JOIN users u ON u.id = j.user_id
		LEFT JOIN job_categories jc ON jc.id = j.job_category_id
		LEFT JOIN job_priorities jp ON jp.id = j.job_priority_id
		LEFT JOIN job_statuses js ON js.id = j.job_status_id
	`
}

// scanDashboardJob escanea una fila en un DashboardJob
func scanDashboardJob(rows *sql.Rows) (domainDashboard.DashboardJob, error) {
	var j domainDashboard.DashboardJob
	var dateReceived time.Time
	var dispatchDate, dueDate, completionDate sql.NullTime
	var workOrder, quickNotes, scheduledTime sql.NullString
	var userID sql.NullInt64
	var weekNumber, routeNumber sql.NullInt64
	var jobSalesPrice sql.NullFloat64
	var propertyStreet, propertyCity, propertyZip, customerName, technicianName sql.NullString
	var categoryName, priorityName, priorityClass, statusName, statusClass sql.NullString

	err := rows.Scan(
		&j.ID, &workOrder, &dateReceived, &j.JobCategoryID, &j.JobPriorityID,
		&j.JobStatusID, &j.PropertyID, &userID, &dispatchDate, &dueDate,
		&quickNotes, &j.Closed, &j.CallAttempted, &j.CageRequired, &weekNumber,
		&routeNumber, &jobSalesPrice, &scheduledTime, &completionDate,
		&propertyStreet, &propertyCity, &propertyZip,
		&customerName,
		&technicianName,
		&categoryName,
		&priorityName, &priorityClass,
		&statusName, &statusClass,
	)
	if err != nil {
		return j, err
	}

	j.DateReceived = dateReceived.Format(timeFormat)

	if workOrder.Valid {
		j.WorkOrder = &workOrder.String
	}
	if userID.Valid {
		j.UserID = &userID.Int64
	}
	if dispatchDate.Valid {
		s := dispatchDate.Time.Format(timeFormat)
		j.DispatchDate = &s
	}
	if dueDate.Valid {
		s := dueDate.Time.Format(timeFormat)
		j.DueDate = &s
	}
	if completionDate.Valid {
		s := completionDate.Time.Format(timeFormat)
		j.CompletionDate = &s
	}
	if quickNotes.Valid {
		j.QuickNotes = &quickNotes.String
	}
	if weekNumber.Valid {
		wn := int(weekNumber.Int64)
		j.WeekNumber = &wn
	}
	if routeNumber.Valid {
		rn := int(routeNumber.Int64)
		j.RouteNumber = &rn
	}
	if jobSalesPrice.Valid {
		j.JobSalesPrice = &jobSalesPrice.Float64
	}
	if scheduledTime.Valid {
		j.ScheduledTime = &scheduledTime.String
	}
	if propertyStreet.Valid {
		j.PropertyStreet = &propertyStreet.String
	}
	if propertyCity.Valid {
		j.PropertyCity = &propertyCity.String
	}
	if propertyZip.Valid {
		j.PropertyZip = &propertyZip.String
	}
	if customerName.Valid {
		j.CustomerName = &customerName.String
	}
	if technicianName.Valid {
		j.TechnicianName = &technicianName.String
	}
	if categoryName.Valid {
		j.CategoryName = &categoryName.String
	}
	if priorityName.Valid {
		j.PriorityName = &priorityName.String
	}
	if priorityClass.Valid {
		j.PriorityClass = &priorityClass.String
	}
	if statusName.Valid {
		j.StatusName = &statusName.String
	}
	if statusClass.Valid {
		j.StatusClass = &statusClass.String
	}

	return j, nil
}

// GetJobs obtiene jobs con filtros, ordenamiento, búsqueda y paginación
func (r *Repository) GetJobs(ctx context.Context, baseFilters map[string]interface{}, filters domainDashboard.DashboardFilters) (*domainDashboard.DashboardJobsResult, error) {
	// Construir query base
	baseQuery := dashboardJobQuery()
	var whereClauses []string
	var args []interface{}

	// Siempre excluir soft deleted
	whereClauses = append(whereClauses, "j.deleted_at IS NULL")

	// Aplicar filtros base (obligatorios del use case)
	if val, ok := baseFilters["job_status_id"].(int64); ok {
		whereClauses = append(whereClauses, "j.job_status_id = ?")
		args = append(args, val)
	}
	if val, ok := baseFilters["job_priority_id"].(int64); ok {
		whereClauses = append(whereClauses, "j.job_priority_id = ?")
		args = append(args, val)
	}
	if val, ok := baseFilters["closed"].(bool); ok {
		whereClauses = append(whereClauses, "j.closed = ?")
		args = append(args, val)
	}
	if val, ok := baseFilters["user_id"].(int64); ok {
		whereClauses = append(whereClauses, "j.user_id = ?")
		args = append(args, val)
	}

	// Aplicar filtros adicionales del usuario
	if filters.UserID != nil {
		if *filters.UserID == -1 {
			// Unassigned
			whereClauses = append(whereClauses, "j.user_id IS NULL")
		} else {
			whereClauses = append(whereClauses, "j.user_id = ?")
			args = append(args, *filters.UserID)
		}
	}

	if filters.JobStatusID != nil {
		whereClauses = append(whereClauses, "j.job_status_id = ?")
		args = append(args, *filters.JobStatusID)
	}

	if filters.JobPriorityID != nil {
		whereClauses = append(whereClauses, "j.job_priority_id = ?")
		args = append(args, *filters.JobPriorityID)
	}

	if filters.IsClosed != nil {
		switch *filters.IsClosed {
		case "all":
			// No agregar filtro
		case "1":
			whereClauses = append(whereClauses, "j.closed = 1")
		default:
			whereClauses = append(whereClauses, "j.closed = 0")
		}
	}

	if filters.Year != nil {
		whereClauses = append(whereClauses, "YEAR(j.date_received) = ?")
		args = append(args, *filters.Year)
	}

	if filters.Week != nil {
		whereClauses = append(whereClauses, "j.week_number = ?")
		args = append(args, *filters.Week)
	}

	if filters.LastDays != nil && *filters.LastDays > 0 {
		whereClauses = append(whereClauses, "j.date_received >= DATE_SUB(NOW(), INTERVAL ? DAY)")
		args = append(args, *filters.LastDays)
	}

	// Búsqueda de texto libre
	if filters.Search != nil && *filters.Search != "" {
		searchPattern := "%" + *filters.Search + "%"
		searchClause := `(
			j.work_order LIKE ? OR
			p.property_code LIKE ? OR
			p.street LIKE ? OR
			p.city LIKE ? OR
			p.state LIKE ? OR
			p.zip LIKE ? OR
			c.name LIKE ? OR
			c.email LIKE ? OR
			c.phone LIKE ? OR
			c.mobile LIKE ?
		)`
		whereClauses = append(whereClauses, searchClause)
		for i := 0; i < 10; i++ {
			args = append(args, searchPattern)
		}
	}

	// Construir WHERE clause
	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = " WHERE " + whereClauses[0]
		for i := 1; i < len(whereClauses); i++ {
			whereSQL += " AND " + whereClauses[i]
		}
	}

	// Contar total de resultados
	countQuery := "SELECT COUNT(*) FROM jobs j LEFT JOIN properties p ON p.id = j.property_id LEFT JOIN customers c ON c.id = p.customer_id" + whereSQL
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		slog.ErrorContext(ctx, "Failed to count jobs", slog.String("error", err.Error()))
		return nil, err
	}

	// Ordenamiento
	orderBy := "ORDER BY j.created_at DESC"
	if filters.Sort != nil && *filters.Sort != "" {
		direction := "DESC"
		if filters.Direction != nil && *filters.Direction == "asc" {
			direction = "ASC"
		}

		switch *filters.Sort {
		case "work_order":
			orderBy = "ORDER BY j.work_order " + direction
		case "date_received", "created_at":
			orderBy = "ORDER BY j.date_received " + direction
		case "dispatch_date":
			orderBy = "ORDER BY j.dispatch_date " + direction
		case "due_date":
			orderBy = "ORDER BY j.due_date " + direction
		case "completion_date":
			orderBy = "ORDER BY j.completion_date " + direction
		case "week_number":
			orderBy = "ORDER BY j.week_number " + direction
		case "job_sales_price":
			orderBy = "ORDER BY j.job_sales_price " + direction
		case "property.city":
			orderBy = "ORDER BY p.city " + direction
		case "property.zip":
			orderBy = "ORDER BY p.zip " + direction
		case "property.customer.name":
			orderBy = "ORDER BY c.name " + direction
		case "user_id":
			orderBy = "ORDER BY u.name " + direction
		case "status":
			orderBy = "ORDER BY js.label " + direction
		case "priority.order":
			orderBy = "ORDER BY jp.order " + direction
		default:
			orderBy = "ORDER BY j.created_at DESC"
		}
	}

	// Paginación
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 20
	}
	offset := (filters.Page - 1) * filters.PageSize

	// Query final
	finalQuery := baseQuery + whereSQL + " " + orderBy + " LIMIT ? OFFSET ?"
	finalArgs := append(args, filters.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, finalQuery, finalArgs...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get jobs", slog.String("error", err.Error()))
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var jobs []domainDashboard.DashboardJob
	for rows.Next() {
		j, err := scanDashboardJob(rows)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to scan dashboard job", slog.String("error", err.Error()))
			return nil, err
		}
		jobs = append(jobs, j)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if jobs == nil {
		jobs = []domainDashboard.DashboardJob{}
	}

	totalPages := (total + filters.PageSize - 1) / filters.PageSize

	return &domainDashboard.DashboardJobsResult{
		Jobs:       jobs,
		Total:      total,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		TotalPages: totalPages,
	}, nil
}
