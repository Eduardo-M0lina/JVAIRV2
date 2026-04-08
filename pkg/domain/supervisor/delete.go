package supervisor

import (
	"context"
	"errors"
	"log/slog"
)

func (uc *UseCase) Delete(ctx context.Context, id int64) error {
	existing, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get supervisor for deletion",
			slog.String("error", err.Error()),
			slog.Int64("supervisor_id", id))
		return err
	}

	if existing.DeletedAt != nil {
		slog.WarnContext(ctx, "Supervisor already deleted",
			slog.Int64("supervisor_id", id))
		return errors.New("supervisor already deleted")
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "Failed to delete supervisor",
			slog.String("error", err.Error()),
			slog.Int64("supervisor_id", id))
		return err
	}

	slog.InfoContext(ctx, "Supervisor deleted successfully",
		slog.Int64("supervisor_id", id))

	return nil
}
