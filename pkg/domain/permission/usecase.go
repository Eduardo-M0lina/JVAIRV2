package permission

import (
	"context"

	"github.com/your-org/jvairv2/pkg/domain/ability"
)

// UseCase define los casos de uso para la gestión de permisos
type UseCase struct {
	repo        Repository
	abilityRepo ability.Repository
}

// NewUseCase crea una nueva instancia del caso de uso de permisos
func NewUseCase(repo Repository, abilityRepo ability.Repository) *UseCase {
	return &UseCase{
		repo:        repo,
		abilityRepo: abilityRepo,
	}
}

// GetByID obtiene un permiso por su ID
func (uc *UseCase) GetByID(ctx context.Context, id int64) (*Permission, error) {
	return uc.repo.GetByID(ctx, id)
}

// GetByEntity obtiene todos los permisos para una entidad específica
func (uc *UseCase) GetByEntity(ctx context.Context, entityType string, entityID int64) ([]*Permission, error) {
	return uc.repo.GetByEntity(ctx, entityType, entityID)
}

// GetByAbility obtiene todos los permisos para una ability específica
func (uc *UseCase) GetByAbility(ctx context.Context, abilityID int64) ([]*Permission, error) {
	return uc.repo.GetByAbility(ctx, abilityID)
}

// Create crea un nuevo permiso
func (uc *UseCase) Create(ctx context.Context, permission *Permission) error {
	// Verificar que la ability exista
	_, err := uc.abilityRepo.GetByID(ctx, permission.AbilityID)
	if err != nil {
		return err
	}

	return uc.repo.Create(ctx, permission)
}

// Update actualiza un permiso existente
func (uc *UseCase) Update(ctx context.Context, permission *Permission) error {
	// Verificar que la ability exista
	_, err := uc.abilityRepo.GetByID(ctx, permission.AbilityID)
	if err != nil {
		return err
	}

	return uc.repo.Update(ctx, permission)
}

// Delete elimina un permiso
func (uc *UseCase) Delete(ctx context.Context, id int64) error {
	return uc.repo.Delete(ctx, id)
}

// Exists verifica si existe un permiso específico
func (uc *UseCase) Exists(ctx context.Context, abilityID, entityID int64, entityType string) (bool, error) {
	return uc.repo.Exists(ctx, abilityID, entityID, entityType)
}

// List obtiene una lista paginada de permisos con filtros opcionales
func (uc *UseCase) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Permission, int, error) {
	return uc.repo.List(ctx, filters, page, pageSize)
}
