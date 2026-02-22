package job_category

import "context"

type Repository interface {
	Create(ctx context.Context, category *JobCategory) error
	GetByID(ctx context.Context, id int64) (*JobCategory, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*JobCategory, int, error)
	Update(ctx context.Context, category *JobCategory) error
	Delete(ctx context.Context, id int64) error
	HasJobs(ctx context.Context, id int64) (bool, error)
}
