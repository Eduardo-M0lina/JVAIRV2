package technician_job_status

import (
	"context"
	"log/slog"
)

func (uc *UseCase) GetByID(ctx context.Context, id int64) (*TechnicianJobStatus, error) {
	status, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get technician job status by ID",
			slog.String("error", err.Error()),
			slog.Int64("technician_job_status_id", id))
		return nil, err
	}

	slog.InfoContext(ctx, "Technician job status retrieved successfully",
		slog.Int64("technician_job_status_id", id))

	return status, nil
}
