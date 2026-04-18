package property

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/angumol/jvairv2/pkg/domain/property"
)

func (r *Repository) SearchByAddress(ctx context.Context, address string) ([]*property.Property, error) {
	if len(address) < 3 {
		return nil, fmt.Errorf("address must be at least 3 characters long")
	}

	query := `
		SELECT
			p.id, p.customer_id, c.name as customer_name, p.property_code, p.street, p.city, p.state, p.zip, p.notes,
			p.created_at, p.updated_at, p.deleted_at
		FROM properties p
		INNER JOIN customers c ON p.customer_id = c.id
		WHERE p.deleted_at IS NULL
		AND p.street LIKE ?
		ORDER BY p.street ASC
	`

	searchPattern := "%" + address + "%"

	rows, err := r.db.QueryContext(ctx, query, searchPattern)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to search properties by address",
			slog.String("error", err.Error()),
			slog.String("address", address))
		return nil, err
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
			&p.CustomerName,
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
			return nil, err
		}
		properties = append(properties, p)
	}

	if err = rows.Err(); err != nil {
		slog.ErrorContext(ctx, "Error iterating property rows",
			slog.String("error", err.Error()))
		return nil, err
	}

	return properties, nil
}
