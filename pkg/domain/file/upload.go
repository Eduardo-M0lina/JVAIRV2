package file

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"
)

// Upload sube un archivo y crea el registro en la base de datos
func (uc *UseCase) Upload(ctx context.Context, fileableID int64, fileableType string, filename string, contentType string, body io.Reader) (*File, error) {
	// Generar key única para S3
	key := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), filename)

	// Subir a S3
	url, err := uc.storage.Upload(ctx, key, body, contentType)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to upload file to storage",
			slog.String("error", err.Error()),
			slog.String("filename", filename))
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Crear registro en la base de datos
	f := &File{
		Type:         &contentType,
		Path:         &key,
		URL:          url,
		FileableID:   fileableID,
		FileableType: fileableType,
	}

	id, err := uc.repo.Create(ctx, f)
	if err != nil {
		// Intentar eliminar el archivo de S3 si falla la creación en BD
		_ = uc.storage.Delete(ctx, key)
		slog.ErrorContext(ctx, "Failed to create file record",
			slog.String("error", err.Error()))
		return nil, err
	}

	f.ID = id
	slog.InfoContext(ctx, "File uploaded successfully",
		slog.Int64("id", id),
		slog.String("key", key))

	return f, nil
}
