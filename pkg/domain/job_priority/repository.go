package job_priority

import "context"

type Repository interface {
	Create(ctx context.Context, priority *JobPriority) error
	GetByID(ctx context.Context, id int64) (*JobPriority, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*JobPriority, int, error)
	Update(ctx context.Context, priority *JobPriority) error
	Delete(ctx context.Context, id int64) error
	HasJobs(ctx context.Context, id int64) (bool, error)
}
