package job_status

import "context"

type Service interface {
	Create(ctx context.Context, status *JobStatus) error
	GetByID(ctx context.Context, id int64) (*JobStatus, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*JobStatus, int, error)
	Update(ctx context.Context, status *JobStatus) error
	Delete(ctx context.Context, id int64) error
}

type UseCase struct {
	repo Repository
}

func NewUseCase(repo Repository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}
