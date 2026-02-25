package supervisor

import (
	"context"
	"log/slog"
)

func (uc *UseCase) GetByID(ctx context.Context, id int64) (*Supervisor, error) {
	supervisor, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get supervisor by ID",
			slog.String("error", err.Error()),
			slog.Int64("supervisor_id", id))
		return nil, err
	}

	slog.InfoContext(ctx, "Supervisor retrieved successfully",
		slog.Int64("supervisor_id", id))

	return supervisor, nil
}
