package quote

import "context"

// JobChecker verifica existencia de jobs
type JobChecker interface {
	GetByID(ctx context.Context, id int64) (interface{}, error)
}

// QuoteStatusChecker verifica existencia de estados de cotización
type QuoteStatusChecker interface {
	GetByID(ctx context.Context, id int64) (interface{}, error)
}

// Service define la interfaz del caso de uso de cotizaciones
type Service interface {
	Create(ctx context.Context, q *Quote) error
	GetByID(ctx context.Context, id int64) (*Quote, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Quote, int64, error)
	Update(ctx context.Context, q *Quote) error
	Delete(ctx context.Context, id int64) error
}

// UseCase implementa la lógica de negocio de cotizaciones
type UseCase struct {
	repo            Repository
	jobRepo         JobChecker
	quoteStatusRepo QuoteStatusChecker
}

// NewUseCase crea una nueva instancia del caso de uso de cotizaciones
func NewUseCase(
	repo Repository,
	jobRepo JobChecker,
	quoteStatusRepo QuoteStatusChecker,
) *UseCase {
	return &UseCase{
		repo:            repo,
		jobRepo:         jobRepo,
		quoteStatusRepo: quoteStatusRepo,
	}
}
