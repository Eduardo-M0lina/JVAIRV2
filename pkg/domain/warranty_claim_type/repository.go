package warranty_claim_type

import "context"

type Repository interface {
	Create(ctx context.Context, wct *WarrantyClaimType) error
	GetByID(ctx context.Context, id int64) (*WarrantyClaimType, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*WarrantyClaimType, int, error)
	Update(ctx context.Context, wct *WarrantyClaimType) error
	Delete(ctx context.Context, id int64) error
	HasClaims(ctx context.Context, id int64) (bool, error)
}
