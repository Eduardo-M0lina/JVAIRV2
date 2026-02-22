package task_status

import "context"

type Repository interface {
	Create(ctx context.Context, status *TaskStatus) error
	GetByID(ctx context.Context, id int64) (*TaskStatus, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*TaskStatus, int, error)
	Update(ctx context.Context, status *TaskStatus) error
	Delete(ctx context.Context, id int64) error
	HasJobTasks(ctx context.Context, id int64) (bool, error)
}
