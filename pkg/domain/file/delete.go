package file

import (
	"context"
	"log/slog"
)

// Delete elimina un archivo de S3 y de la base de datos
func (uc *UseCase) Delete(ctx context.Context, id int64) error {
	// Obtener el archivo para saber el path en S3
	f, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if err == ErrFileNotFound {
			return err
		}
		slog.ErrorContext(ctx, "Failed to get file for delete",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return err
	}

	// Eliminar de S3
	if f.Path != nil && *f.Path != "" {
		if err := uc.storage.Delete(ctx, *f.Path); err != nil {
			slog.ErrorContext(ctx, "Failed to delete file from storage",
				slog.String("error", err.Error()),
				slog.String("path", *f.Path))
			return err
		}
	}

	// Eliminar registro de la base de datos
	if err := uc.repo.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "Failed to delete file record",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return err
	}

	slog.InfoContext(ctx, "File deleted successfully",
		slog.Int64("id", id))

	return nil
}
