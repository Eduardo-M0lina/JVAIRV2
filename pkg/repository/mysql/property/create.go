package property

import (
	"context"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/property"
)

func (r *Repository) Create(ctx context.Context, p *property.Property) error {
	query := `
		INSERT INTO properties (
			customer_id, property_code, street, city, state, zip, notes,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		p.CustomerID,
		p.PropertyCode,
		p.Street,
		p.City,
		p.State,
		p.Zip,
		p.Notes,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute insert property query",
			slog.String("error", err.Error()))
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get last insert ID",
			slog.String("error", err.Error()))
		return err
	}

	p.ID = id
	return nil
}
