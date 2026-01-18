package property

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/property"
)

func (r *Repository) GetByID(ctx context.Context, id int64) (*property.Property, error) {
	query := `
		SELECT
			id, customer_id, property_code, street, city, state, zip, notes,
			created_at, updated_at, deleted_at
		FROM properties
		WHERE id = ? AND deleted_at IS NULL
	`

	p := &property.Property{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
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
