package warranty_status

import "context"

type Service interface {
	Create(ctx context.Context, ws *WarrantyStatus) error
	GetByID(ctx context.Context, id int64) (*WarrantyStatus, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*WarrantyStatus, int, error)
	Update(ctx context.Context, ws *WarrantyStatus) error
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
