package warranty_claim

import "context"

// Repository define los métodos para interactuar con el almacenamiento de warranty claims
type Repository interface {
	Create(ctx context.Context, wc *WarrantyClaim) error
	GetByID(ctx context.Context, id int64) (*WarrantyClaim, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*WarrantyClaim, int, error)
	Update(ctx context.Context, wc *WarrantyClaim) error
	Delete(ctx context.Context, id int64) error
}
