package task_status

import (
	"context"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/task_status"
)

func (r *Repository) Update(ctx context.Context, s *task_status.TaskStatus) error {
	query := "UPDATE task_statuses SET label = ?, class = ?, `order` = ?, is_active = ?, updated_at = NOW() WHERE id = ?"

	_, err := r.db.ExecContext(ctx, query,
		s.Label,
		s.Class,
		s.Order,
		s.IsActive,
		s.ID,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to update task status",
			slog.String("error", err.Error()),
			slog.Int64("id", s.ID))
		return err
	}

	return nil
}
