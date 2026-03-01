package job_visit

import (
	"context"
	"log/slog"

	domainJV "github.com/your-org/jvairv2/pkg/domain/job_visit"
)

// Create crea una nueva visita de trabajo en la base de datos
func (r *Repository) Create(ctx context.Context, visit *domainJV.JobVisit) (int64, error) {
	query := `
		INSERT INTO job_visits (job_id, user_id, viewable_by, date, report, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		visit.JobID,
		visit.UserID,
		visit.ViewableBy,
		visit.Date,
		visit.Report,
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create job visit",
			slog.String("error", err.Error()))
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get last insert ID for job visit",
			slog.String("error", err.Error()))
		return 0, err
	}

	return id, nil
}
