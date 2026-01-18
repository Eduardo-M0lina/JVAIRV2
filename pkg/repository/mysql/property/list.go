package property

import (
	"context"
	"log/slog"
	"strings"

	"github.com/your-org/jvairv2/pkg/domain/property"
)

func (r *Repository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*property.Property, int, error) {
	baseQuery := `
		SELECT
			id, customer_id, property_code, street, city, state, zip, notes,
			created_at, updated_at, deleted_at
		FROM properties
		WHERE deleted_at IS NULL
	`

	countQuery := "SELECT COUNT(*) FROM properties WHERE deleted_at IS NULL"

	var args []interface{}
	var conditions []string

	if customerID, ok := filters["customer_id"].(int64); ok && customerID > 0 {
		conditions = append(conditions, "customer_id = ?")
		args = append(args, customerID)
	}

	if search, ok := filters["search"].(string); ok && search != "" {
		searchCondition := `(
			property_code LIKE ? OR
			street LIKE ? OR
			city LIKE ? OR
			state LIKE ? OR
			zip LIKE ?
		)`
		conditions = append(conditions, searchCondition)
		searchPattern := "%" + search + "%"
		for i := 0; i < 5; i++ {
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
		slog.ErrorContext(ctx, "Failed to count properties",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	baseQuery += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	offset := (page - 1) * pageSize
	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to query properties",
			slog.String("error", err.Error()))
		return nil, 0, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			slog.ErrorContext(ctx, "Failed to close rows", slog.String("error", closeErr.Error()))
		}
	}()

	var properties []*property.Property
	for rows.Next() {
		p := &property.Property{}
		err := rows.Scan(
			&p.ID,
			&p.CustomerID,
			&p.PropertyCode,
			&p.Street,
			&p.City,
			&p.State,
			&p.Zip,
			&p.Notes,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.DeletedAt,
		)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to scan property row",
				slog.String("error", err.Error()))
			return nil, 0, err
		}
		properties = append(properties, p)
	}

	if err = rows.Err(); err != nil {
		slog.ErrorContext(ctx, "Error iterating property rows",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	return properties, total, nil
}
