package role

import (
	"context"
)

// UseCase define los casos de uso para la gestión de roles
type UseCase struct {
	repo Repository
}

// NewUseCase crea una nueva instancia del caso de uso de roles
func NewUseCase(repo Repository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

// GetByID obtiene un rol por su ID
func (uc *UseCase) GetByID(ctx context.Context, id int64) (*Role, error) {
	return uc.repo.GetByID(ctx, id)
}

// GetByName obtiene un rol por su nombre
func (uc *UseCase) GetByName(ctx context.Context, name string) (*Role, error) {
	return uc.repo.GetByName(ctx, name)
}

// Create crea un nuevo rol
func (uc *UseCase) Create(ctx context.Context, role *Role) error {
	return uc.repo.Create(ctx, role)
}

// Update actualiza un rol existente
func (uc *UseCase) Update(ctx context.Context, role *Role) error {
	return uc.repo.Update(ctx, role)
}

// Delete elimina un rol
func (uc *UseCase) Delete(ctx context.Context, id int64) error {
	return uc.repo.Delete(ctx, id)
}

// List obtiene una lista paginada de roles con filtros opcionales
func (uc *UseCase) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Role, int, error) {
	return uc.repo.List(ctx, filters, page, pageSize)
}
