package file

import (
	"context"
	"log/slog"
)

// GetByID obtiene un archivo por su ID
func (uc *UseCase) GetByID(ctx context.Context, id int64) (*File, error) {
	f, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if err == ErrFileNotFound {
			return nil, err
		}
		slog.ErrorContext(ctx, "Failed to get file",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return nil, err
	}
	return f, nil
}
