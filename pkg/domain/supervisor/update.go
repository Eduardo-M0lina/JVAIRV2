package supervisor

import (
	"context"
	"errors"
	"log/slog"
)

func (uc *UseCase) Update(ctx context.Context, supervisor *Supervisor) error {
	// Validar que el supervisor existe
	existing, err := uc.repo.GetByID(ctx, supervisor.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get supervisor for update",
			slog.String("error", err.Error()),
			slog.Int64("supervisor_id", supervisor.ID))
		return err
	}

	// Validar que no esté eliminado
	if existing.DeletedAt != nil {
		slog.WarnContext(ctx, "Cannot update deleted supervisor",
			slog.Int64("supervisor_id", supervisor.ID))
		return errors.New("cannot update deleted supervisor")
	}

	// Validar que el customer existe
	_, err = uc.customerRepo.GetByID(ctx, supervisor.CustomerID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to validate customer",
			slog.String("error", err.Error()),
			slog.Int64("customer_id", supervisor.CustomerID))
		return errors.New("invalid customer ID")
	}

	if err := uc.repo.Update(ctx, supervisor); err != nil {
		slog.ErrorContext(ctx, "Failed to update supervisor",
			slog.String("error", err.Error()),
			slog.Int64("supervisor_id", supervisor.ID))
		return err
	}

	slog.InfoContext(ctx, "Supervisor updated successfully",
		slog.Int64("supervisor_id", supervisor.ID),
		slog.String("name", supervisor.Name))

	return nil
}
