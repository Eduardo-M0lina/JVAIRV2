package job_visit

import "context"

// Repository define las operaciones de persistencia para visitas de trabajo
type Repository interface {
	Create(ctx context.Context, visit *JobVisit) (int64, error)
	GetByID(ctx context.Context, id int64) (*JobVisit, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int, sort, direction string) ([]*JobVisit, int64, error)
	Update(ctx context.Context, visit *JobVisit) error
	Delete(ctx context.Context, id int64) error
}

// JobChecker verifica la existencia de un job
type JobChecker interface {
	JobExists(ctx context.Context, id int64) (bool, error)
}

// UserChecker verifica la existencia de un usuario
type UserChecker interface {
	UserExists(ctx context.Context, id int64) (bool, error)
}
