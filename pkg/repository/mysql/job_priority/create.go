package job_priority

import (
	"context"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/job_priority"
)

func (r *Repository) Create(ctx context.Context, p *job_priority.JobPriority) error {
	query := `
		INSERT INTO job_priorities (label, ` + "`order`" + `, class, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		p.Label,
		p.Order,
		p.Class,
		p.IsActive,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute insert job priority query",
			slog.String("error", err.Error()))
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get last insert ID",
			slog.String("error", err.Error()))
		return err
	}

	p.ID = id
	return nil
}
