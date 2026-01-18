package customer

import (
	"context"
	"errors"
	"log/slog"
)

func (uc *UseCase) Create(ctx context.Context, customer *Customer) error {
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

	if err := uc.repo.Create(ctx, customer); err != nil {
		slog.ErrorContext(ctx, "Failed to create customer",
			slog.String("error", err.Error()),
			slog.String("name", customer.Name))
		return err
	}

	slog.InfoContext(ctx, "Customer created successfully",
		slog.Int64("customer_id", customer.ID),
		slog.String("name", customer.Name))

	return nil
}
