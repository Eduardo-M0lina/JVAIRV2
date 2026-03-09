package job_resident

import "context"

type Repository interface {
	Create(ctx context.Context, resident *JobResident) error
	GetByID(ctx context.Context, id int64) (*JobResident, error)
	List(ctx context.Context, jobID int64, limit, offset int) ([]*JobResident, int64, error)
	Update(ctx context.Context, resident *JobResident) error
	Delete(ctx context.Context, id int64) error
}

type JobExistsChecker interface {
	JobExists(ctx context.Context, jobID int64) (bool, error)
}
