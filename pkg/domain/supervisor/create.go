package supervisor

import (
	"context"
	"errors"
	"log/slog"
)

func (uc *UseCase) Create(ctx context.Context, supervisor *Supervisor) error {
	// Validar que el customer existe
	_, err := uc.customerRepo.GetByID(ctx, supervisor.CustomerID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to validate customer",
			slog.String("error", err.Error()),
			slog.Int64("customer_id", supervisor.CustomerID))
		return errors.New("invalid customer ID")
	}

	if err := uc.repo.Create(ctx, supervisor); err != nil {
		slog.ErrorContext(ctx, "Failed to create supervisor",
			slog.String("error", err.Error()),
			slog.String("name", supervisor.Name))
		return err
	}

	slog.InfoContext(ctx, "Supervisor created successfully",
		slog.Int64("supervisor_id", supervisor.ID),
		slog.String("name", supervisor.Name))

	return nil
}
