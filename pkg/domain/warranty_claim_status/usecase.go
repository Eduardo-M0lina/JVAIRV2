package warranty_claim_status

import "context"

type Service interface {
	Create(ctx context.Context, wcs *WarrantyClaimStatus) error
	GetByID(ctx context.Context, id int64) (*WarrantyClaimStatus, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*WarrantyClaimStatus, int, error)
	Update(ctx context.Context, wcs *WarrantyClaimStatus) error
	Delete(ctx context.Context, id int64) error
}

type UseCase struct {
	repo Repository
}

func NewUseCase(repo Repository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}
