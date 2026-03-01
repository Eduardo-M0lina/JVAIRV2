package job_visit

import (
	"context"
	"log/slog"
)

// List obtiene una lista paginada de visitas de trabajo
func (uc *UseCase) List(ctx context.Context, filters map[string]interface{}, page, pageSize int, sort, direction string) ([]*JobVisit, int64, error) {
	visits, total, err := uc.repo.List(ctx, filters, page, pageSize, sort, direction)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list job visits",
			slog.String("error", err.Error()))
		return nil, 0, err
	}
	return visits, total, nil
}
