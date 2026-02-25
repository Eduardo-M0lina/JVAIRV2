package supervisor

import (
	"context"
	"log/slog"
	"strings"

	"github.com/your-org/jvairv2/pkg/domain/supervisor"
)

func (r *Repository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*supervisor.Supervisor, int, error) {
	baseQuery := `
		SELECT
			id, customer_id, name, phone, email, created_at, updated_at, deleted_at
		FROM supervisors
		WHERE deleted_at IS NULL
	`

	countQuery := "SELECT COUNT(*) FROM supervisors WHERE deleted_at IS NULL"

	var args []interface{}
	var conditions []string

	if customerID, ok := filters["customer_id"].(int64); ok && customerID > 0 {
		conditions = append(conditions, "customer_id = ?")
		args = append(args, customerID)
	}

	if search, ok := filters["search"].(string); ok && search != "" {
		searchCondition := `(
			name LIKE ? OR
			phone LIKE ? OR
			email LIKE ?
		)`
		conditions = append(conditions, searchCondition)
		searchPattern := "%" + search + "%"
		for i := 0; i < 3; i++ {
			args = append(args, searchPattern)
		}
	}

	if len(conditions) > 0 {
		whereClause := " AND " + strings.Join(conditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
	}

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to count supervisors",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	baseQuery += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	offset := (page - 1) * pageSize
	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to query supervisors",
			slog.String("error", err.Error()))
		return nil, 0, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			slog.ErrorContext(ctx, "Failed to close rows", slog.String("error", closeErr.Error()))
		}
	}()

	var supervisors []*supervisor.Supervisor
	for rows.Next() {
		s := &supervisor.Supervisor{}
		err := rows.Scan(
			&s.ID,
			&s.CustomerID,
			&s.Name,
			&s.Phone,
			&s.Email,
			&s.CreatedAt,
			&s.UpdatedAt,
			&s.DeletedAt,
		)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to scan supervisor row",
				slog.String("error", err.Error()))
			return nil, 0, err
		}
		supervisors = append(supervisors, s)
	}

	if err = rows.Err(); err != nil {
		slog.ErrorContext(ctx, "Error iterating supervisor rows",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	return supervisors, total, nil
}
