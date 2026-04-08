package job

import (
	"context"
	"log/slog"
)

// GetByID obtiene un job por su ID
func (uc *UseCase) GetByID(ctx context.Context, id int64) (*Job, error) {
	j, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get job",
			slog.Int64("id", id),
			slog.String("error", err.Error()))
		return nil, err
	}

	if j.IsDeleted() {
		return nil, ErrJobNotFound
	}

	return j, nil
}
