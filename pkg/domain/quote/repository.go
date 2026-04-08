package quote

import "context"

type Repository interface {
	Create(ctx context.Context, q *Quote) error
	GetByID(ctx context.Context, id int64) (*Quote, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Quote, int64, error)
	Update(ctx context.Context, q *Quote) error
	Delete(ctx context.Context, id int64) error
}
