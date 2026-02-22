package job_category

import (
	"context"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/job_category"
)

func (r *Repository) Update(ctx context.Context, c *job_category.JobCategory) error {
	query := `
		UPDATE job_categories
		SET label = ?, label_plural = ?, type = ?, is_active = ?, updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		c.Label,
		c.LabelPlural,
		c.Type,
		c.IsActive,
		c.ID,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to update job category",
			slog.String("error", err.Error()),
			slog.Int64("id", c.ID))
		return err
	}

	return nil
}
