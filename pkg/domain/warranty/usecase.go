package warranty

import "context"

// Service define la interfaz del servicio de warranties
type Service interface {
	Create(ctx context.Context, w *Warranty) error
	GetByID(ctx context.Context, id int64) (*Warranty, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Warranty, int, error)
	Update(ctx context.Context, w *Warranty) error
	Delete(ctx context.Context, id int64) error
}

// JobChecker verifica existencia de jobs
type JobChecker interface {
	GetByID(ctx context.Context, id int64) (interface{}, error)
}

// WarrantyTypeChecker verifica existencia de tipos de garantía
type WarrantyTypeChecker interface {
	GetByID(ctx context.Context, id int64) (interface{}, error)
}

// WarrantyStatusChecker verifica existencia de estados de garantía
type WarrantyStatusChecker interface {
	GetByID(ctx context.Context, id int64) (interface{}, error)
}

// UseCase implementa la lógica de negocio de warranties
type UseCase struct {
	repo        Repository
	jobCheck    JobChecker
	typeCheck   WarrantyTypeChecker
	statusCheck WarrantyStatusChecker
}

// NewUseCase crea una nueva instancia del caso de uso de warranties
func NewUseCase(repo Repository, jobCheck JobChecker, typeCheck WarrantyTypeChecker, statusCheck WarrantyStatusChecker) *UseCase {
	return &UseCase{
		repo:        repo,
		jobCheck:    jobCheck,
		typeCheck:   typeCheck,
		statusCheck: statusCheck,
	}
}
