package warranty_type

import "context"

type Repository interface {
	Create(ctx context.Context, wt *WarrantyType) error
	GetByID(ctx context.Context, id int64) (*WarrantyType, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*WarrantyType, int, error)
	Update(ctx context.Context, wt *WarrantyType) error
	Delete(ctx context.Context, id int64) error
	HasWarranties(ctx context.Context, id int64) (bool, error)
}
