package job_priority

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/your-org/jvairv2/pkg/domain/job_priority"
)

func (r *Repository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*job_priority.JobPriority, int, error) {
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
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM job_priorities WHERE %s", whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		slog.ErrorContext(ctx, "Failed to count job priorities",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	// Query with pagination
	offset := (page - 1) * pageSize
	query := fmt.Sprintf("SELECT id, label, `order`, class, is_active, created_at, updated_at FROM job_priorities WHERE %s ORDER BY `order` ASC LIMIT ? OFFSET ?", whereClause)

	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list job priorities",
			slog.String("error", err.Error()))
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var priorities []*job_priority.JobPriority
	for rows.Next() {
		p := &job_priority.JobPriority{}
		if err := rows.Scan(
			&p.ID,
			&p.Label,
			&p.Order,
			&p.Class,
			&p.IsActive,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			slog.ErrorContext(ctx, "Failed to scan job priority",
				slog.String("error", err.Error()))
			return nil, 0, err
		}
		priorities = append(priorities, p)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return priorities, total, nil
}
