package quote

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Delete(ctx context.Context, id int64) error {
	// Verificar que la cotización existe
	_, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get quote for deletion",
			slog.String("error", err.Error()),
			slog.Int64("quote_id", id))
		return err
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "Failed to delete quote",
			slog.String("error", err.Error()),
			slog.Int64("quote_id", id))
		return err
	}

	slog.InfoContext(ctx, "Quote deleted successfully",
		slog.Int64("quote_id", id))

	return nil
}
