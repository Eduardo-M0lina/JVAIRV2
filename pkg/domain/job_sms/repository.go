package job_sms

import "context"

type Repository interface {
	Create(ctx context.Context, item *JobSMS) error
	GetByID(ctx context.Context, id int64) (*JobSMS, error)
	List(ctx context.Context, jobID int64, limit, offset int) ([]*JobSMS, int64, error)
	Delete(ctx context.Context, id int64) error
}

type JobExistsChecker interface {
	JobExists(ctx context.Context, jobID int64) (bool, error)
}
