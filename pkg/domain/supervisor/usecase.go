package supervisor

import (
	"context"

	"github.com/your-org/jvairv2/pkg/domain/customer"
)

type Service interface {
	Create(ctx context.Context, supervisor *Supervisor) error
	GetByID(ctx context.Context, id int64) (*Supervisor, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Supervisor, int, error)
	Update(ctx context.Context, supervisor *Supervisor) error
	Delete(ctx context.Context, id int64) error
}

type UseCase struct {
	repo         Repository
	customerRepo customer.Repository
}

func NewUseCase(repo Repository, customerRepo customer.Repository) *UseCase {
	return &UseCase{
		repo:         repo,
		customerRepo: customerRepo,
	}
}
