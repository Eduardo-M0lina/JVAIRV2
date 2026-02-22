package job_status

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/job_status"
)

func (r *Repository) GetByID(ctx context.Context, id int64) (*job_status.JobStatus, error) {
	query := `SELECT id, label, class, is_active, created_at, updated_at FROM job_statuses WHERE id = ?`

	s := &job_status.JobStatus{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID,
		&s.Label,
		&s.Class,
		&s.IsActive,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, job_status.ErrJobStatusNotFound
		}
		slog.ErrorContext(ctx, "Failed to get job status by ID",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return nil, err
	}

	return s, nil
}
