package technician_job_status

import (
	"context"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/technician_job_status"
)

func (r *Repository) Update(ctx context.Context, s *technician_job_status.TechnicianJobStatus) error {
	query := `
		UPDATE technician_job_statuses
		SET label = ?, class = ?, job_status_id = ?, is_active = ?, updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		s.Label,
		s.Class,
		s.JobStatusID,
		s.IsActive,
		s.ID,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to update technician job status",
			slog.String("error", err.Error()),
			slog.Int64("id", s.ID))
		return err
	}

	return nil
}
