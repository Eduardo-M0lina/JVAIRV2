package warranty_type

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/your-org/jvairv2/pkg/domain/warranty_type"
)

func (r *Repository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*warranty_type.WarrantyType, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}

	if search, ok := filters["search"].(string); ok && search != "" {
		where = append(where, "(label LIKE ? OR label_plural LIKE ?)")
		args = append(args, "%"+search+"%", "%"+search+"%")
	}

	whereClause := strings.Join(where, " AND ")

	// Count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM warranty_types WHERE %s", whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		slog.ErrorContext(ctx, "Failed to count warranty types",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	// Query
	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT id, label, label_plural, is_active, created_at, updated_at
		FROM warranty_types
		WHERE %s
		ORDER BY label ASC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list warranty types",
			slog.String("error", err.Error()))
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var types []*warranty_type.WarrantyType
	for rows.Next() {
		wt := &warranty_type.WarrantyType{}

		if err := rows.Scan(
			&wt.ID,
			&wt.Label,
			&wt.LabelPlural,
			&wt.IsActive,
			&wt.CreatedAt,
			&wt.UpdatedAt,
		); err != nil {
			slog.ErrorContext(ctx, "Failed to scan warranty type row",
				slog.String("error", err.Error()))
			return nil, 0, err
		}

		types = append(types, wt)
	}

	return types, total, nil
}
