package property

import (
	"context"
	"log/slog"
)

func (uc *UseCase) GetByID(ctx context.Context, id int64) (*Property, error) {
	property, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get property by ID",
			slog.String("error", err.Error()),
			slog.Int64("property_id", id))
		return nil, err
	}

	slog.InfoContext(ctx, "Property retrieved successfully",
		slog.Int64("property_id", id))

	return property, nil
}
