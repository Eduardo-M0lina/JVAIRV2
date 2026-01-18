package property

import (
	"github.com/your-org/jvairv2/pkg/domain/customer"
)

// UseCase orquesta las operaciones de negocio para propiedades
type UseCase struct {
	repo         Repository
	customerRepo customer.Repository
}

// NewUseCase crea una nueva instancia de UseCase
func NewUseCase(repo Repository, customerRepo customer.Repository) *UseCase {
	return &UseCase{
		repo:         repo,
		customerRepo: customerRepo,
	}
}
