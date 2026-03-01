package file

import (
	"context"
	"log/slog"
)

// ListByFileable obtiene los archivos asociados a una entidad
func (uc *UseCase) ListByFileable(ctx context.Context, fileableID int64, fileableType string) ([]*File, error) {
	files, err := uc.repo.ListByFileable(ctx, fileableID, fileableType)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list files",
			slog.String("error", err.Error()),
			slog.Int64("fileableId", fileableID),
			slog.String("fileableType", fileableType))
		return nil, err
	}
	return files, nil
}
