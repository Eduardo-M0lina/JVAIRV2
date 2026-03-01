package warranty_status

import "context"

type Repository interface {
	Create(ctx context.Context, ws *WarrantyStatus) error
	GetByID(ctx context.Context, id int64) (*WarrantyStatus, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*WarrantyStatus, int, error)
	Update(ctx context.Context, ws *WarrantyStatus) error
	Delete(ctx context.Context, id int64) error
	HasWarranties(ctx context.Context, id int64) (bool, error)
}
