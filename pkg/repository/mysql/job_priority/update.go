package job_priority

import (
	"context"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/job_priority"
)

func (r *Repository) Update(ctx context.Context, p *job_priority.JobPriority) error {
	query := "UPDATE job_priorities SET label = ?, `order` = ?, class = ?, is_active = ?, updated_at = NOW() WHERE id = ?"

	_, err := r.db.ExecContext(ctx, query,
		p.Label,
		p.Order,
		p.Class,
		p.IsActive,
		p.ID,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to update job priority",
			slog.String("error", err.Error()),
			slog.Int64("id", p.ID))
		return err
	}

	return nil
}
