package supervisor

import "context"

type Repository interface {
	Create(ctx context.Context, supervisor *Supervisor) error
	GetByID(ctx context.Context, id int64) (*Supervisor, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Supervisor, int, error)
	Update(ctx context.Context, supervisor *Supervisor) error
	Delete(ctx context.Context, id int64) error
}
