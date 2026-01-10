package ability

import (
	"context"
)

// UseCase define los casos de uso para la gestión de abilities (capacidades/permisos)
type UseCase struct {
	repo Repository
}

// NewUseCase crea una nueva instancia del caso de uso de abilities
func NewUseCase(repo Repository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

// GetByID obtiene una ability por su ID
func (uc *UseCase) GetByID(ctx context.Context, id int64) (*Ability, error) {
	return uc.repo.GetByID(ctx, id)
}

// GetByName obtiene una ability por su nombre
func (uc *UseCase) GetByName(ctx context.Context, name string) (*Ability, error) {
	return uc.repo.GetByName(ctx, name)
}

// Create crea una nueva ability
func (uc *UseCase) Create(ctx context.Context, ability *Ability) error {
	return uc.repo.Create(ctx, ability)
}

// Update actualiza una ability existente
func (uc *UseCase) Update(ctx context.Context, ability *Ability) error {
	return uc.repo.Update(ctx, ability)
}

// Delete elimina una ability
func (uc *UseCase) Delete(ctx context.Context, id int64) error {
	return uc.repo.Delete(ctx, id)
}

// List obtiene una lista paginada de abilities con filtros opcionales
func (uc *UseCase) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Ability, int, error) {
	return uc.repo.List(ctx, filters, page, pageSize)
}
