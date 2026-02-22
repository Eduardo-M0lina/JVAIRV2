package quote_status

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Delete(ctx context.Context, id int64) error {
	// Validar que el status existe
	_, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get quote status for deletion",
			slog.String("error", err.Error()),
			slog.Int64("quote_status_id", id))
		return err
	}

	// Verificar que no tenga quotes asociadas
	hasQuotes, err := uc.repo.HasQuotes(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to check quote status quotes",
			slog.String("error", err.Error()),
			slog.Int64("quote_status_id", id))
		return err
	}

	if hasQuotes {
		slog.WarnContext(ctx, "Cannot delete quote status with quotes",
			slog.Int64("quote_status_id", id))
		return ErrQuoteStatusInUse
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "Failed to delete quote status",
			slog.String("error", err.Error()),
			slog.Int64("quote_status_id", id))
		return err
	}

	slog.InfoContext(ctx, "Quote status deleted successfully",
		slog.Int64("quote_status_id", id))

	return nil
}
