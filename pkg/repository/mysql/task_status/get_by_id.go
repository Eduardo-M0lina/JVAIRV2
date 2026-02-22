package task_status

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/task_status"
)

func (r *Repository) GetByID(ctx context.Context, id int64) (*task_status.TaskStatus, error) {
	query := "SELECT id, label, class, `order`, is_active, created_at, updated_at FROM task_statuses WHERE id = ?"

	s := &task_status.TaskStatus{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID,
		&s.Label,
		&s.Class,
		&s.Order,
		&s.IsActive,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, task_status.ErrTaskStatusNotFound
		}
		slog.ErrorContext(ctx, "Failed to get task status by ID",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return nil, err
	}

	return s, nil
}
