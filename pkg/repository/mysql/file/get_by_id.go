package file

import (
	"context"
	"database/sql"
	"log/slog"

	domainFile "github.com/angumol/jvairv2/pkg/domain/file"
)

// GetByID obtiene un archivo por su ID
func (r *Repository) GetByID(ctx context.Context, id int64) (*domainFile.File, error) {
	query := `
		SELECT id, type, path, url, fileable_id, fileable_type, created_at, updated_at
		FROM files
		WHERE id = ?
	`

	f := &domainFile.File{}
	var fileType, path sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&f.ID,
		&fileType,
		&path,
		&f.URL,
		&f.FileableID,
		&f.FileableType,
		&f.CreatedAt,
		&f.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domainFile.ErrFileNotFound
		}
		slog.ErrorContext(ctx, "Failed to get file by ID",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return nil, err
	}

	if fileType.Valid {
		f.Type = &fileType.String
	}
	if path.Valid {
		f.Path = &path.String
	}

	return f, nil
}
