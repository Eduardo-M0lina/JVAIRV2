package property

import (
	"context"
	"errors"
	"log/slog"
)

func (uc *UseCase) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Property, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Validar que el customer existe si se proporciona como filtro
	if customerID, ok := filters["customer_id"].(int64); ok && customerID > 0 {
		customer, err := uc.customerRepo.GetByID(ctx, customerID)
		if err != nil {
			slog.WarnContext(ctx, "Invalid customer_id in filters",
				slog.Int64("customer_id", customerID),
				slog.String("error", err.Error()))
			return nil, 0, errors.New("invalid customer_id")
		}
		if customer.DeletedAt != nil {
			slog.WarnContext(ctx, "Customer is deleted",
				slog.Int64("customer_id", customerID))
			return nil, 0, errors.New("customer is deleted")
		}
	}

	properties, total, err := uc.repo.List(ctx, filters, page, pageSize)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list properties",
			slog.String("error", err.Error()),
			slog.Int("page", page),
			slog.Int("pageSize", pageSize))
		return nil, 0, err
	}

	slog.InfoContext(ctx, "Properties listed successfully",
		slog.Int("total", total),
		slog.Int("page", page),
		slog.Int("pageSize", pageSize))

	return properties, total, nil
}
