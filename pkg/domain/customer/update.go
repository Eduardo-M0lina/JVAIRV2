package customer

import (
	"context"
	"errors"
	"log/slog"
)

func (uc *UseCase) Update(ctx context.Context, customer *Customer) error {
	// Validar que el cliente existe
	existing, err := uc.repo.GetByID(ctx, customer.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get customer for update",
			slog.String("error", err.Error()),
			slog.Int64("customer_id", customer.ID))
		return err
	}

	// Validar que no esté eliminado (regla de negocio)
	if existing.DeletedAt != nil {
		slog.WarnContext(ctx, "Cannot update deleted customer",
			slog.Int64("customer_id", customer.ID))
		return errors.New("cannot update deleted customer")
	}

	// Validar que el workflow existe y está activo (regla de negocio)
	workflow, err := uc.workflowRepo.GetByID(ctx, customer.WorkflowID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to validate workflow",
			slog.String("error", err.Error()),
			slog.Int64("workflow_id", customer.WorkflowID))
		return errors.New("invalid workflow ID")
	}

	if !workflow.IsActive {
		slog.WarnContext(ctx, "Workflow is not active",
			slog.Int64("workflow_id", customer.WorkflowID))
		return errors.New("workflow is not active")
	}

	if err := uc.repo.Update(ctx, customer); err != nil {
		slog.ErrorContext(ctx, "Failed to update customer",
			slog.String("error", err.Error()),
			slog.Int64("customer_id", customer.ID))
		return err
	}

	slog.InfoContext(ctx, "Customer updated successfully",
		slog.Int64("customer_id", customer.ID),
		slog.String("name", customer.Name))

	return nil
}
