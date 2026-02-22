package technician_job_status

import (
	"context"
	"errors"
	"log/slog"
)

func (uc *UseCase) Update(ctx context.Context, status *TechnicianJobStatus) error {
	if err := status.Validate(); err != nil {
		return err
	}

	// Validar que el estado existe
	_, err := uc.repo.GetByID(ctx, status.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get technician job status for update",
			slog.String("error", err.Error()),
			slog.Int64("technician_job_status_id", status.ID))
		return err
	}

	// Validar que el job_status_id existe si se proporciona
	if status.JobStatusID != nil {
		_, err := uc.jobStatusRepo.GetByID(ctx, *status.JobStatusID)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to validate job status",
				slog.String("error", err.Error()),
				slog.Int64("job_status_id", *status.JobStatusID))
			return errors.New("invalid job_status_id")
		}
	}

	if err := uc.repo.Update(ctx, status); err != nil {
		slog.ErrorContext(ctx, "Failed to update technician job status",
			slog.String("error", err.Error()),
			slog.Int64("technician_job_status_id", status.ID))
		return err
	}

	slog.InfoContext(ctx, "Technician job status updated successfully",
		slog.Int64("technician_job_status_id", status.ID),
		slog.String("label", status.Label))

	return nil
}
