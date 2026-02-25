package supervisor

import (
	"context"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/supervisor"
)

func (r *Repository) Update(ctx context.Context, s *supervisor.Supervisor) error {
	query := `
		UPDATE supervisors SET
			customer_id = ?,
			name = ?,
			phone = ?,
			email = ?,
			updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		s.CustomerID,
		s.Name,
		s.Phone,
		s.Email,
		s.ID,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to update supervisor",
			slog.String("error", err.Error()),
			slog.Int64("supervisor_id", s.ID))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get rows affected",
			slog.String("error", err.Error()))
		return err
	}

	if rowsAffected == 0 {
		slog.WarnContext(ctx, "No rows affected during update",
			slog.Int64("supervisor_id", s.ID))
	}

	return nil
}
