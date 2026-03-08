package job_activity_log

import "context"

type Repository interface {
	Create(ctx context.Context, log *JobActivityLog) error
	GetByID(ctx context.Context, id int64) (*JobActivityLog, error)
	List(ctx context.Context, jobID int64, limit, offset int) ([]*JobActivityLog, int64, error)
	Delete(ctx context.Context, id int64) error
}

type JobExistsChecker interface {
	JobExists(ctx context.Context, jobID int64) (bool, error)
}

type UserExistsChecker interface {
	UserExists(ctx context.Context, userID int64) (bool, error)
}
