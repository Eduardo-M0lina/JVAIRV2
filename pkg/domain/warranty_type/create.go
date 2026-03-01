package warranty_type

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Create(ctx context.Context, wt *WarrantyType) error {
	if err := wt.Validate(); err != nil {
		return err
	}

	if err := uc.repo.Create(ctx, wt); err != nil {
		slog.ErrorContext(ctx, "Failed to create warranty type",
			slog.String("error", err.Error()),
			slog.String("label", wt.Label))
		return err
	}

	slog.InfoContext(ctx, "Warranty type created successfully",
		slog.Int64("warranty_type_id", wt.ID),
		slog.String("label", wt.Label))

	return nil
}
