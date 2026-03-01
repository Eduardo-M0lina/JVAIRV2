package warranty

import (
	"context"
	"log/slog"
)

// Delete elimina una garantía (soft delete)
func (uc *UseCase) Delete(ctx context.Context, id int64) error {
	existing, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get warranty for deletion",
			slog.String("error", err.Error()),
			slog.Int64("warranty_id", id))
		return err
	}

	if existing.IsDeleted() {
		return ErrWarrantyDeleted
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "Failed to delete warranty",
			slog.String("error", err.Error()),
			slog.Int64("warranty_id", id))
		return err
	}

	slog.InfoContext(ctx, "Warranty deleted successfully",
		slog.Int64("warranty_id", id))

	return nil
}
