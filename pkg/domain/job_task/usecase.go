package job_task

import "context"

type Service interface {
	Create(ctx context.Context, task *JobTask) error
	GetByID(ctx context.Context, id int64) (*JobTask, error)
	List(ctx context.Context, jobID int64, limit, offset int) ([]*JobTask, int64, error)
	ListAll(ctx context.Context, limit, offset int) ([]*JobTask, int64, error)
	Update(ctx context.Context, task *JobTask) error
	Delete(ctx context.Context, id int64) error
}

type service struct {
	repo              Repository
	jobChecker        JobExistsChecker
	userChecker       UserExistsChecker
	taskStatusChecker TaskStatusExistsChecker
}

func NewUseCase(repo Repository, jobChecker JobExistsChecker, userChecker UserExistsChecker, taskStatusChecker TaskStatusExistsChecker) Service {
	return &service{
		repo:              repo,
		jobChecker:        jobChecker,
		userChecker:       userChecker,
		taskStatusChecker: taskStatusChecker,
	}
}
