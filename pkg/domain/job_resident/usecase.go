package job_resident

import "context"

type Service interface {
	Create(ctx context.Context, resident *JobResident) error
	GetByID(ctx context.Context, id int64) (*JobResident, error)
	List(ctx context.Context, jobID int64, limit, offset int) ([]*JobResident, int64, error)
	Update(ctx context.Context, resident *JobResident) error
	Delete(ctx context.Context, id int64) error
}

type service struct {
	repo       Repository
	jobChecker JobExistsChecker
}

func NewUseCase(repo Repository, jobChecker JobExistsChecker) Service {
	return &service{
		repo:       repo,
		jobChecker: jobChecker,
	}
}
