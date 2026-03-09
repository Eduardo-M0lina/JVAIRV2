package job_rate_status

import "context"

type Repository interface {
	Create(ctx context.Context, status *JobRateStatus) error
	GetByID(ctx context.Context, id int64) (*JobRateStatus, error)
	List(ctx context.Context) ([]*JobRateStatus, error)
	Update(ctx context.Context, status *JobRateStatus) error
	Delete(ctx context.Context, id int64) error
}
