package file

import (
	"context"
	"log/slog"

	domainFile "github.com/your-org/jvairv2/pkg/domain/file"
)

// Create crea un nuevo registro de archivo en la base de datos
func (r *Repository) Create(ctx context.Context, f *domainFile.File) (int64, error) {
	query := `
		INSERT INTO files (type, path, url, fileable_id, fileable_type, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		f.Type,
		f.Path,
		f.URL,
		f.FileableID,
		f.FileableType,
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create file record",
			slog.String("error", err.Error()))
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get last insert ID for file",
			slog.String("error", err.Error()))
		return 0, err
	}

	return id, nil
}
