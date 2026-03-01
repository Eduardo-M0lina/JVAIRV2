package file

import (
	"context"
	"io"
	"log/slog"
	"path/filepath"
)

// Download descarga un archivo de S3
func (uc *UseCase) Download(ctx context.Context, id int64) (io.ReadCloser, string, string, error) {
	f, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if err == ErrFileNotFound {
			return nil, "", "", err
		}
		slog.ErrorContext(ctx, "Failed to get file for download",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return nil, "", "", err
	}

	if f.Path == nil || *f.Path == "" {
		return nil, "", "", ErrFileNotFound
	}

	body, contentType, err := uc.storage.GetObject(ctx, *f.Path)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to download file from storage",
			slog.String("error", err.Error()),
			slog.String("path", *f.Path))
		return nil, "", "", err
	}

	filename := filepath.Base(*f.Path)

	return body, contentType, filename, nil
}
