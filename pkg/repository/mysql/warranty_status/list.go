package warranty_status

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/angumol/jvairv2/pkg/domain/warranty_status"
)

func (r *Repository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*warranty_status.WarrantyStatus, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}

	if search, ok := filters["search"].(string); ok && search != "" {
		where = append(where, "label LIKE ?")
		args = append(args, "%"+search+"%")
	}

	whereClause := strings.Join(where, " AND ")

	// Count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM warranty_statuses WHERE %s", whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		slog.ErrorContext(ctx, "Failed to count warranty statuses",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	// Query
	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT id, label, class, `+"`order`"+`, is_active, created_at, updated_at
		FROM warranty_statuses
		WHERE %s
		ORDER BY `+"`order`"+` ASC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list warranty statuses",
			slog.String("error", err.Error()))
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var statuses []*warranty_status.WarrantyStatus
	for rows.Next() {
		ws := &warranty_status.WarrantyStatus{}
		var class sql.NullString

		if err := rows.Scan(
			&ws.ID,
			&ws.Label,
			&class,
			&ws.Order,
			&ws.IsActive,
			&ws.CreatedAt,
			&ws.UpdatedAt,
		); err != nil {
			slog.ErrorContext(ctx, "Failed to scan warranty status row",
				slog.String("error", err.Error()))
			return nil, 0, err
		}

		if class.Valid {
			ws.Class = &class.String
		}

		statuses = append(statuses, ws)
	}

	return statuses, total, nil
}
