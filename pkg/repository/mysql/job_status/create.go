package job_status

import (
	"context"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/job_status"
)

func (r *Repository) Create(ctx context.Context, s *job_status.JobStatus) error {
	query := `
		INSERT INTO job_statuses (label, class, is_active, created_at, updated_at)
		VALUES (?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		s.Label,
		s.Class,
		s.IsActive,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute insert job status query",
			slog.String("error", err.Error()))
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get last insert ID",
			slog.String("error", err.Error()))
		return err
	}

	s.ID = id
	return nil
}
