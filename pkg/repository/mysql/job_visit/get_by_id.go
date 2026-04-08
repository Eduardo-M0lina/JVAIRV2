package job_visit

import (
	"context"
	"database/sql"
	"log/slog"

	domainJV "github.com/angumol/jvairv2/pkg/domain/job_visit"
)

// GetByID obtiene una visita de trabajo por su ID
func (r *Repository) GetByID(ctx context.Context, id int64) (*domainJV.JobVisit, error) {
	query := `
		SELECT id, job_id, user_id, viewable_by, date, report,
			created_at, updated_at, deleted_at
		FROM job_visits
		WHERE id = ? AND deleted_at IS NULL
	`

	jv := &domainJV.JobVisit{}
	var viewableBy, report sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&jv.ID,
		&jv.JobID,
		&jv.UserID,
		&viewableBy,
		&jv.Date,
		&report,
		&jv.CreatedAt,
		&jv.UpdatedAt,
		&jv.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domainJV.ErrJobVisitNotFound
		}
		slog.ErrorContext(ctx, "Failed to get job visit by ID",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return nil, err
	}

	if viewableBy.Valid {
		jv.ViewableBy = &viewableBy.String
	}
	if report.Valid {
		jv.Report = &report.String
	}

	return jv, nil
}
