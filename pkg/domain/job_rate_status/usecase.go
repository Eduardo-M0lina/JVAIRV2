package job_rate_status

import "context"

type Service interface {
	Create(ctx context.Context, status *JobRateStatus) error
	GetByID(ctx context.Context, id int64) (*JobRateStatus, error)
	List(ctx context.Context) ([]*JobRateStatus, error)
	Update(ctx context.Context, status *JobRateStatus) error
	Delete(ctx context.Context, id int64) error
}

type service struct {
	repo Repository
}

func NewUseCase(repo Repository) Service {
	return &service{repo: repo}
}
