package warranty

import (
	"context"
	"log/slog"
)

// GetByID obtiene una garantía por su ID
func (uc *UseCase) GetByID(ctx context.Context, id int64) (*Warranty, error) {
	w, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get warranty by ID",
			slog.String("error", err.Error()),
			slog.Int64("warranty_id", id))
		return nil, err
	}

	slog.InfoContext(ctx, "Warranty retrieved successfully",
		slog.Int64("warranty_id", id))

	return w, nil
}
