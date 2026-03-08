package job_task

import "context"

type Repository interface {
	Create(ctx context.Context, task *JobTask) error
	GetByID(ctx context.Context, id int64) (*JobTask, error)
	List(ctx context.Context, jobID int64, limit, offset int) ([]*JobTask, int64, error)
	ListAll(ctx context.Context, limit, offset int) ([]*JobTask, int64, error)
	Update(ctx context.Context, task *JobTask) error
	Delete(ctx context.Context, id int64) error
}

type JobExistsChecker interface {
	JobExists(ctx context.Context, jobID int64) (bool, error)
}

type UserExistsChecker interface {
	UserExists(ctx context.Context, userID int64) (bool, error)
}

type TaskStatusExistsChecker interface {
	TaskStatusExists(ctx context.Context, taskStatusID int64) (bool, error)
}
