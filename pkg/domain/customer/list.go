package customer

import (
	"context"
	"errors"
	"log/slog"
)

func (uc *UseCase) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Customer, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Validar que el workflow existe si se proporciona como filtro
	if workflowID, ok := filters["workflow_id"].(int64); ok && workflowID > 0 {
		workflow, err := uc.workflowRepo.GetByID(ctx, workflowID)
		if err != nil {
			slog.WarnContext(ctx, "Invalid workflow_id in filters",
				slog.Int64("workflow_id", workflowID),
				slog.String("error", err.Error()))
			return nil, 0, errors.New("invalid workflow_id")
		}
		if !workflow.IsActive {
			slog.WarnContext(ctx, "Workflow is not active",
				slog.Int64("workflow_id", workflowID))
			return nil, 0, errors.New("workflow is not active")
		}
	}

	customers, total, err := uc.repo.List(ctx, filters, page, pageSize)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list customers",
			slog.String("error", err.Error()),
			slog.Int("page", page),
			slog.Int("page_size", pageSize))
		return nil, 0, err
	}

	slog.InfoContext(ctx, "Customers listed successfully",
		slog.Int("total", total),
		slog.Int("page", page),
		slog.Int("page_size", pageSize))

	return customers, total, nil
}
