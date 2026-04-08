package customer

import "context"

type Repository interface {
	Create(ctx context.Context, customer *Customer) error
	GetByID(ctx context.Context, id int64) (*Customer, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Customer, int, error)
	Update(ctx context.Context, customer *Customer) error
	Delete(ctx context.Context, id int64) error
	HasProperties(ctx context.Context, id int64) (bool, error)
}
