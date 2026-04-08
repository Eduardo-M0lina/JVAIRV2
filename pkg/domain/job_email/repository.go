package job_email

import "context"

type Repository interface {
	Create(ctx context.Context, item *JobEmail) error
	GetByID(ctx context.Context, id int64) (*JobEmail, error)
	List(ctx context.Context, jobID int64, limit, offset int) ([]*JobEmail, int64, error)
	Delete(ctx context.Context, id int64) error
}

type JobExistsChecker interface {
	JobExists(ctx context.Context, jobID int64) (bool, error)
}
