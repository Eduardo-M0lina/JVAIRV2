package job_status

import (
	"context"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/job_status"
)

func (r *Repository) Update(ctx context.Context, s *job_status.JobStatus) error {
	query := `
		UPDATE job_statuses
		SET label = ?, class = ?, is_active = ?, updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		s.Label,
		s.Class,
		s.IsActive,
		s.ID,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to update job status",
			slog.String("error", err.Error()),
			slog.Int64("id", s.ID))
		return err
	}

	return nil
}
