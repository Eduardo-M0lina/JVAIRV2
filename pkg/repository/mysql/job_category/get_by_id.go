package job_category

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/job_category"
)

func (r *Repository) GetByID(ctx context.Context, id int64) (*job_category.JobCategory, error) {
	query := `SELECT id, label, label_plural, type, is_active, created_at, updated_at FROM job_categories WHERE id = ?`

	category := &job_category.JobCategory{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&category.ID,
		&category.Label,
		&category.LabelPlural,
		&category.Type,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, job_category.ErrJobCategoryNotFound
		}
		slog.ErrorContext(ctx, "Failed to get job category by ID",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return nil, err
	}

	return category, nil
}
