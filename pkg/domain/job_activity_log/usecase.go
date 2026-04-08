package job_activity_log

import "context"

type Service interface {
	Create(ctx context.Context, log *JobActivityLog) error
	GetByID(ctx context.Context, id int64) (*JobActivityLog, error)
	List(ctx context.Context, jobID int64, limit, offset int) ([]*JobActivityLog, int64, error)
	Delete(ctx context.Context, id int64) error
}

type service struct {
	repo        Repository
	jobChecker  JobExistsChecker
	userChecker UserExistsChecker
}

func NewUseCase(repo Repository, jobChecker JobExistsChecker, userChecker UserExistsChecker) Service {
	return &service{
		repo:        repo,
		jobChecker:  jobChecker,
		userChecker: userChecker,
	}
}
