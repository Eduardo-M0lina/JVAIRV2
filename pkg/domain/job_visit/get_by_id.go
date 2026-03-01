package job_visit

import (
	"context"
	"log/slog"
)

// GetByID obtiene una visita de trabajo por su ID
func (uc *UseCase) GetByID(ctx context.Context, id int64) (*JobVisit, error) {
	visit, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if err == ErrJobVisitNotFound {
			return nil, err
		}
		slog.ErrorContext(ctx, "Failed to get job visit",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return nil, err
	}
	return visit, nil
}
