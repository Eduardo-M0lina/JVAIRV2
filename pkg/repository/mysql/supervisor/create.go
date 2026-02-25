package supervisor

import (
	"context"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/supervisor"
)

func (r *Repository) Create(ctx context.Context, s *supervisor.Supervisor) error {
	query := `
		INSERT INTO supervisors (
			customer_id, name, phone, email, created_at, updated_at
		) VALUES (?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		s.CustomerID,
		s.Name,
		s.Phone,
		s.Email,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute insert supervisor query",
			slog.String("error", err.Error()))
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get last insert ID",
			slog.String("error", err.Error()))
		return err
	}

	s.ID = id
	return nil
}
