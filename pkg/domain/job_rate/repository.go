package job_rate

import "context"

type Repository interface {
	Create(ctx context.Context, rate *JobRate) error
	GetByID(ctx context.Context, id int64) (*JobRate, error)
	List(ctx context.Context, jobID int64, limit, offset int) ([]*JobRate, int64, error)
	Update(ctx context.Context, rate *JobRate) error
	Delete(ctx context.Context, id int64) error
}

type JobExistsChecker interface {
	JobExists(ctx context.Context, jobID int64) (bool, error)
}

type UserExistsChecker interface {
	UserExists(ctx context.Context, userID int64) (bool, error)
}

type JobRateStatusExistsChecker interface {
	JobRateStatusExists(ctx context.Context, statusID int64) (bool, error)
}
