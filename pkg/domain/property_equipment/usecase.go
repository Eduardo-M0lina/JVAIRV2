package property_equipment

import (
	"github.com/angumol/jvairv2/pkg/domain/property"
)

// UseCase orquesta las operaciones de negocio para equipos de propiedad
type UseCase struct {
	repo         Repository
	propertyRepo property.Repository
}

// NewUseCase crea una nueva instancia de UseCase
func NewUseCase(repo Repository, propertyRepo property.Repository) *UseCase {
	return &UseCase{
		repo:         repo,
		propertyRepo: propertyRepo,
	}
}
