package job_priority

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/job_priority"
)

func (r *Repository) GetByID(ctx context.Context, id int64) (*job_priority.JobPriority, error) {
	query := "SELECT id, label, `order`, class, is_active, created_at, updated_at FROM job_priorities WHERE id = ?"

	p := &job_priority.JobPriority{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.Label,
		&p.Order,
		&p.Class,
		&p.IsActive,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, job_priority.ErrJobPriorityNotFound
		}
		slog.ErrorContext(ctx, "Failed to get job priority by ID",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return nil, err
	}

	return p, nil
}
