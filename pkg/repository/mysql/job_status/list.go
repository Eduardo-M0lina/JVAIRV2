package job_status

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/your-org/jvairv2/pkg/domain/job_status"
)

func (r *Repository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*job_status.JobStatus, int, error) {
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
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM job_statuses WHERE %s", whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		slog.ErrorContext(ctx, "Failed to count job statuses",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	// Query with pagination
	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT id, label, class, is_active, created_at, updated_at
		FROM job_statuses
		WHERE %s
		ORDER BY id ASC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list job statuses",
			slog.String("error", err.Error()))
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var statuses []*job_status.JobStatus
	for rows.Next() {
		s := &job_status.JobStatus{}
		if err := rows.Scan(
			&s.ID,
			&s.Label,
			&s.Class,
			&s.IsActive,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			slog.ErrorContext(ctx, "Failed to scan job status",
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
