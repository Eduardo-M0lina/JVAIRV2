package file

import (
	"context"
	"database/sql"
	"log/slog"

	domainFile "github.com/your-org/jvairv2/pkg/domain/file"
)

// ListByFileable obtiene los archivos asociados a una entidad polimórfica
func (r *Repository) ListByFileable(ctx context.Context, fileableID int64, fileableType string) ([]*domainFile.File, error) {
	query := `
		SELECT id, type, path, url, fileable_id, fileable_type, created_at, updated_at
		FROM files
		WHERE fileable_id = ? AND fileable_type = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, fileableID, fileableType)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list files by fileable",
			slog.String("error", err.Error()),
			slog.Int64("fileableId", fileableID),
			slog.String("fileableType", fileableType))
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var files []*domainFile.File
	for rows.Next() {
		f := &domainFile.File{}
		var fileType, path sql.NullString

		if err := rows.Scan(
			&f.ID,
			&fileType,
			&path,
			&f.URL,
			&f.FileableID,
			&f.FileableType,
			&f.CreatedAt,
			&f.UpdatedAt,
		); err != nil {
			slog.ErrorContext(ctx, "Failed to scan file row",
				slog.String("error", err.Error()))
			return nil, err
		}

		if fileType.Valid {
			f.Type = &fileType.String
		}
		if path.Valid {
			f.Path = &path.String
		}

		files = append(files, f)
	}

	return files, nil
}
