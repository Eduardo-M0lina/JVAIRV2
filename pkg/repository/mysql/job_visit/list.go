package job_visit

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	domainJV "github.com/your-org/jvairv2/pkg/domain/job_visit"
)

// List obtiene una lista paginada de visitas de trabajo
func (r *Repository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int, sort, direction string) ([]*domainJV.JobVisit, int64, error) {
	baseQuery := `FROM job_visits jv WHERE jv.deleted_at IS NULL`
	args := make([]interface{}, 0)

	// Filtros
	if jobID, ok := filters["jobId"]; ok {
		baseQuery += ` AND jv.job_id = ?`
		args = append(args, jobID)
	}
	if userID, ok := filters["userId"]; ok {
		baseQuery += ` AND jv.user_id = ?`
		args = append(args, userID)
	}
	if search, ok := filters["search"]; ok {
		baseQuery += ` AND (jv.report LIKE ?)`
		args = append(args, fmt.Sprintf("%%%s%%", search))
	}

	// Contar total
	var total int64
	countQuery := `SELECT COUNT(*) ` + baseQuery
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to count job visits",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	// Ordenamiento
	allowedSorts := map[string]string{
		"date":       "jv.date",
		"created_at": "jv.created_at",
	}
	sortColumn, ok := allowedSorts[sort]
	if !ok {
		sortColumn = "jv.created_at"
	}
	if direction != "ASC" && direction != "asc" {
		direction = "DESC"
	}

	// Query principal
	selectQuery := fmt.Sprintf(`SELECT jv.id, jv.job_id, jv.user_id, jv.viewable_by, jv.date, jv.report,
		jv.created_at, jv.updated_at, jv.deleted_at %s ORDER BY %s %s LIMIT ? OFFSET ?`,
		baseQuery, sortColumn, direction)

	offset := (page - 1) * pageSize
	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list job visits",
			slog.String("error", err.Error()))
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var visits []*domainJV.JobVisit
	for rows.Next() {
		jv := &domainJV.JobVisit{}
		var viewableBy, report sql.NullString

		if err := rows.Scan(
			&jv.ID,
			&jv.JobID,
			&jv.UserID,
			&viewableBy,
			&jv.Date,
			&report,
			&jv.CreatedAt,
			&jv.UpdatedAt,
			&jv.DeletedAt,
		); err != nil {
			slog.ErrorContext(ctx, "Failed to scan job visit row",
				slog.String("error", err.Error()))
			return nil, 0, err
		}

		if viewableBy.Valid {
			jv.ViewableBy = &viewableBy.String
		}
		if report.Valid {
			jv.Report = &report.String
		}

		visits = append(visits, jv)
	}

	return visits, total, nil
}
