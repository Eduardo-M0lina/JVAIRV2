package job_category

import (
	"context"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/job_category"
)

func (r *Repository) Create(ctx context.Context, c *job_category.JobCategory) error {
	query := `
		INSERT INTO job_categories (label, label_plural, type, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		c.Label,
		c.LabelPlural,
		c.Type,
		c.IsActive,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute insert job category query",
			slog.String("error", err.Error()))
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get last insert ID",
			slog.String("error", err.Error()))
		return err
	}

	c.ID = id
	return nil
}
