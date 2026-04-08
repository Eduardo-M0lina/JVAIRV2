package invoice

import "context"

// Service define la interfaz del servicio de invoices
type Service interface {
	Create(ctx context.Context, inv *Invoice) error
	GetByID(ctx context.Context, id int64) (*Invoice, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Invoice, int, error)
	Update(ctx context.Context, inv *Invoice) error
	Delete(ctx context.Context, id int64) error
}

// JobChecker verifica existencia de jobs
type JobChecker interface {
	GetByID(ctx context.Context, id int64) (interface{}, error)
}

// UseCase implementa la lógica de negocio de invoices
type UseCase struct {
	repo     Repository
	jobCheck JobChecker
}

// NewUseCase crea una nueva instancia del caso de uso de invoices
func NewUseCase(repo Repository, jobCheck JobChecker) *UseCase {
	return &UseCase{
		repo:     repo,
		jobCheck: jobCheck,
	}
}
