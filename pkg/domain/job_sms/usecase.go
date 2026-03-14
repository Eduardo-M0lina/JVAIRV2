package job_sms

import "context"

type Service interface {
	Create(ctx context.Context, item *JobSMS) error
	GetByID(ctx context.Context, id int64) (*JobSMS, error)
	List(ctx context.Context, jobID int64, limit, offset int) ([]*JobSMS, int64, error)
	Delete(ctx context.Context, id int64) error
}

type service struct {
	repo       Repository
	jobChecker JobExistsChecker
}

func NewUseCase(repo Repository, jobChecker JobExistsChecker) Service {
	return &service{repo: repo, jobChecker: jobChecker}
}
