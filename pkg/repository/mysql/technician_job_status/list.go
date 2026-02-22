package technician_job_status

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/your-org/jvairv2/pkg/domain/technician_job_status"
)

func (r *Repository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*technician_job_status.TechnicianJobStatus, int, error) {
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

	if jobStatusID, ok := filters["job_status_id"].(int64); ok {
		where = append(where, "job_status_id = ?")
		args = append(args, jobStatusID)
	}

	whereClause := strings.Join(where, " AND ")

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM technician_job_statuses WHERE %s", whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		slog.ErrorContext(ctx, "Failed to count technician job statuses",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	// Query with pagination
	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT id, label, class, job_status_id, is_active, created_at, updated_at
		FROM technician_job_statuses
		WHERE %s
		ORDER BY id ASC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list technician job statuses",
			slog.String("error", err.Error()))
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var statuses []*technician_job_status.TechnicianJobStatus
	for rows.Next() {
		s := &technician_job_status.TechnicianJobStatus{}
		if err := rows.Scan(
			&s.ID,
			&s.Label,
			&s.Class,
			&s.JobStatusID,
			&s.IsActive,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			slog.ErrorContext(ctx, "Failed to scan technician job status",
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
