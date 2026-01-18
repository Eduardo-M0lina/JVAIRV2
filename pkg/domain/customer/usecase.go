package customer

import (
	"context"

	"github.com/your-org/jvairv2/pkg/domain/workflow"
)

type Service interface {
	Create(ctx context.Context, customer *Customer) error
	GetByID(ctx context.Context, id int64) (*Customer, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Customer, int, error)
	Update(ctx context.Context, customer *Customer) error
	Delete(ctx context.Context, id int64) error
}

type UseCase struct {
	repo         Repository
	workflowRepo workflow.Repository
}

func NewUseCase(repo Repository, workflowRepo workflow.Repository) *UseCase {
	return &UseCase{
		repo:         repo,
		workflowRepo: workflowRepo,
	}
}
