package warranty_claim

import "context"

// Service define la interfaz del servicio de warranty claims
type Service interface {
	Create(ctx context.Context, wc *WarrantyClaim) error
	GetByID(ctx context.Context, id int64) (*WarrantyClaim, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*WarrantyClaim, int, error)
	Update(ctx context.Context, wc *WarrantyClaim) error
	Delete(ctx context.Context, id int64) error
}

// JobChecker verifica existencia de jobs
type JobChecker interface {
	GetByID(ctx context.Context, id int64) (interface{}, error)
}

// ClaimTypeChecker verifica existencia de tipos de reclamo
type ClaimTypeChecker interface {
	GetByID(ctx context.Context, id int64) (interface{}, error)
}

// ClaimStatusChecker verifica existencia de estados de reclamo
type ClaimStatusChecker interface {
	GetByID(ctx context.Context, id int64) (interface{}, error)
}

// UseCase implementa la lógica de negocio de warranty claims
type UseCase struct {
	repo        Repository
	jobCheck    JobChecker
	typeCheck   ClaimTypeChecker
	statusCheck ClaimStatusChecker
}

// NewUseCase crea una nueva instancia del caso de uso
func NewUseCase(repo Repository, jobCheck JobChecker, typeCheck ClaimTypeChecker, statusCheck ClaimStatusChecker) *UseCase {
	return &UseCase{
		repo:        repo,
		jobCheck:    jobCheck,
		typeCheck:   typeCheck,
		statusCheck: statusCheck,
	}
}
