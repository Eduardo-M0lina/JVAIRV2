package job_category

import "context"

type Service interface {
	Create(ctx context.Context, category *JobCategory) error
	GetByID(ctx context.Context, id int64) (*JobCategory, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*JobCategory, int, error)
	Update(ctx context.Context, category *JobCategory) error
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
