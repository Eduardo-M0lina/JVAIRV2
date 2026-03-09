package job_rate

import "context"

type Service interface {
	Create(ctx context.Context, rate *JobRate) error
	GetByID(ctx context.Context, id int64) (*JobRate, error)
	List(ctx context.Context, jobID int64, limit, offset int) ([]*JobRate, int64, error)
	Update(ctx context.Context, rate *JobRate) error
	Delete(ctx context.Context, id int64) error
}

type service struct {
	repo          Repository
	jobChecker    JobExistsChecker
	userChecker   UserExistsChecker
	statusChecker JobRateStatusExistsChecker
}

func NewUseCase(repo Repository, jobChecker JobExistsChecker, userChecker UserExistsChecker, statusChecker JobRateStatusExistsChecker) Service {
	return &service{
		repo:          repo,
		jobChecker:    jobChecker,
		userChecker:   userChecker,
		statusChecker: statusChecker,
	}
}
