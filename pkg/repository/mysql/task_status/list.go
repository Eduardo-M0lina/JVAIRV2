package task_status

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/your-org/jvairv2/pkg/domain/task_status"
)

func (r *Repository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*task_status.TaskStatus, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}

	if search, ok := filters["search"].(string); ok && search != "" {
		where = append(where, "(label LIKE ? OR class LIKE ?)")
		args = append(args, "%"+search+"%", "%"+search+"%")
	}

	if isActive, ok := filters["is_active"].(bool); ok {
		where = append(where, "is_active = ?")
		args = append(args, isActive)
	}

	whereClause := strings.Join(where, " AND ")

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM task_statuses WHERE %s", whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		slog.ErrorContext(ctx, "Failed to count task statuses",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	// Query with pagination
	offset := (page - 1) * pageSize
	query := fmt.Sprintf("SELECT id, label, class, `order`, is_active, created_at, updated_at FROM task_statuses WHERE %s ORDER BY `order` ASC LIMIT ? OFFSET ?", whereClause)

	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list task statuses",
			slog.String("error", err.Error()))
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var statuses []*task_status.TaskStatus
	for rows.Next() {
		s := &task_status.TaskStatus{}
		if err := rows.Scan(
			&s.ID,
			&s.Label,
			&s.Class,
			&s.Order,
			&s.IsActive,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			slog.ErrorContext(ctx, "Failed to scan task status",
				slog.String("error", err.Error()))
			return nil, 0, err
		}
		statuses = append(statuses, s)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return statuses, total, nil
}
