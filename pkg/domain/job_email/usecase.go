package job_email

import "context"

type Service interface {
	Create(ctx context.Context, item *JobEmail) error
	GetByID(ctx context.Context, id int64) (*JobEmail, error)
	List(ctx context.Context, jobID int64, limit, offset int) ([]*JobEmail, int64, error)
	Delete(ctx context.Context, id int64) error
}

type service struct {
	repo       Repository
	jobChecker JobExistsChecker
}

func NewUseCase(repo Repository, jobChecker JobExistsChecker) Service {
	return &service{repo: repo, jobChecker: jobChecker}
}
