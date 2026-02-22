package technician_job_status

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/technician_job_status"
)

func (r *Repository) GetByID(ctx context.Context, id int64) (*technician_job_status.TechnicianJobStatus, error) {
	query := `SELECT id, label, class, job_status_id, is_active, created_at, updated_at FROM technician_job_statuses WHERE id = ?`

	s := &technician_job_status.TechnicianJobStatus{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID,
		&s.Label,
		&s.Class,
		&s.JobStatusID,
		&s.IsActive,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, technician_job_status.ErrTechnicianJobStatusNotFound
		}
		slog.ErrorContext(ctx, "Failed to get technician job status by ID",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return nil, err
	}

	return s, nil
}
