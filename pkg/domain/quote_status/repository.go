package quote_status

import "context"

type Repository interface {
	Create(ctx context.Context, qs *QuoteStatus) error
	GetByID(ctx context.Context, id int64) (*QuoteStatus, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*QuoteStatus, int, error)
	Update(ctx context.Context, qs *QuoteStatus) error
	Delete(ctx context.Context, id int64) error
	HasQuotes(ctx context.Context, id int64) (bool, error)
}
