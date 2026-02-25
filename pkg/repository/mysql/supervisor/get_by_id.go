package supervisor

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/supervisor"
)

func (r *Repository) GetByID(ctx context.Context, id int64) (*supervisor.Supervisor, error) {
	query := `
		SELECT
			id, customer_id, name, phone, email, created_at, updated_at, deleted_at
		FROM supervisors
		WHERE id = ? AND deleted_at IS NULL
	`

	s := &supervisor.Supervisor{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
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
