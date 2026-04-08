package job_visit

import (
	"context"
	"log/slog"

	domainJV "github.com/angumol/jvairv2/pkg/domain/job_visit"
)

// Update actualiza una visita de trabajo en la base de datos
func (r *Repository) Update(ctx context.Context, visit *domainJV.JobVisit) error {
	query := `
		UPDATE job_visits
		SET user_id = ?, viewable_by = ?, date = ?, report = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query,
		visit.UserID,
		visit.ViewableBy,
		visit.Date,
		visit.Report,
		visit.ID,
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to update job visit",
			slog.String("error", err.Error()),
			slog.Int64("id", visit.ID))
		return err
	}

	return nil
}
