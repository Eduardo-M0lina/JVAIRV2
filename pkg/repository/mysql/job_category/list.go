package job_category

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/your-org/jvairv2/pkg/domain/job_category"
)

func (r *Repository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*job_category.JobCategory, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}

	if search, ok := filters["search"].(string); ok && search != "" {
		where = append(where, "(label LIKE ? OR label_plural LIKE ?)")
		args = append(args, "%"+search+"%", "%"+search+"%")
	}

	if isActive, ok := filters["is_active"].(bool); ok {
		where = append(where, "is_active = ?")
		args = append(args, isActive)
	}

	if typ, ok := filters["type"].(string); ok && typ != "" {
		where = append(where, "type = ?")
		args = append(args, typ)
	}

	whereClause := strings.Join(where, " AND ")

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM job_categories WHERE %s", whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		slog.ErrorContext(ctx, "Failed to count job categories",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	// Query with pagination
	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT id, label, label_plural, type, is_active, created_at, updated_at
		FROM job_categories
		WHERE %s
		ORDER BY id ASC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list job categories",
			slog.String("error", err.Error()))
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var categories []*job_category.JobCategory
	for rows.Next() {
		c := &job_category.JobCategory{}
		if err := rows.Scan(
			&c.ID,
			&c.Label,
			&c.LabelPlural,
			&c.Type,
			&c.IsActive,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			slog.ErrorContext(ctx, "Failed to scan job category",
				slog.String("error", err.Error()))
			return nil, 0, err
		}
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}
