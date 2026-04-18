package property

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/angumol/jvairv2/pkg/domain/property"
)

func (r *Repository) GetByID(ctx context.Context, id int64) (*property.Property, error) {
	query := `
		SELECT
			p.id, p.customer_id, c.name as customer_name, p.property_code, p.street, p.city, p.state, p.zip, p.notes,
			p.created_at, p.updated_at, p.deleted_at
		FROM properties p
		INNER JOIN customers c ON p.customer_id = c.id
		WHERE p.id = ? AND p.deleted_at IS NULL
	`

	p := &property.Property{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
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
		if errors.Is(err, sql.ErrNoRows) {
			slog.WarnContext(ctx, "Property not found",
				slog.Int64("property_id", id))
			return nil, errors.New("property not found")
		}
		slog.ErrorContext(ctx, "Failed to query property by ID",
			slog.String("error", err.Error()),
			slog.Int64("property_id", id))
		return nil, err
	}

	return p, nil
}
