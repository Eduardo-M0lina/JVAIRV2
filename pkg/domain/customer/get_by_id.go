package customer

import (
	"context"
	"log/slog"
)

func (uc *UseCase) GetByID(ctx context.Context, id int64) (*Customer, error) {
	customer, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get customer by ID",
			slog.String("error", err.Error()),
			slog.Int64("customer_id", id))
		return nil, err
	}

	slog.InfoContext(ctx, "Customer retrieved successfully",
		slog.Int64("customer_id", id))

	return customer, nil
}
