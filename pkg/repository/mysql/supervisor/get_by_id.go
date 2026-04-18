package supervisor

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/angumol/jvairv2/pkg/domain/supervisor"
)

func (r *Repository) GetByID(ctx context.Context, id int64) (*supervisor.Supervisor, error) {
	query := `
		SELECT
			s.id, s.customer_id, c.name as customer_name, s.name, s.phone, s.email,
			s.created_at, s.updated_at, s.deleted_at
		FROM supervisors s
		INNER JOIN customers c ON s.customer_id = c.id
		WHERE s.id = ? AND s.deleted_at IS NULL
	`

	s := &supervisor.Supervisor{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID,
		&s.CustomerID,
		&s.CustomerName,
		&s.Name,
		&s.Phone,
		&s.Email,
		&s.CreatedAt,
		&s.UpdatedAt,
		&s.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.WarnContext(ctx, "Supervisor not found",
				slog.Int64("supervisor_id", id))
			return nil, errors.New("supervisor not found")
		}
		slog.ErrorContext(ctx, "Failed to get supervisor by ID",
			slog.String("error", err.Error()),
			slog.Int64("supervisor_id", id))
		return nil, err
	}

	return s, nil
}
