package assigned_role

import (
	"context"

	"github.com/your-org/jvairv2/pkg/domain/role"
)

// UseCase define los casos de uso para la gestión de asignaciones de roles
type UseCase struct {
	repo     Repository
	roleRepo role.Repository
}

// NewUseCase crea una nueva instancia del caso de uso de asignaciones de roles
func NewUseCase(repo Repository, roleRepo role.Repository) *UseCase {
	return &UseCase{
		repo:     repo,
		roleRepo: roleRepo,
	}
}

// GetByID obtiene una asignación de rol por su ID
func (uc *UseCase) GetByID(ctx context.Context, id int64) (*AssignedRole, error) {
	return uc.repo.GetByID(ctx, id)
}

// GetByEntity obtiene todas las asignaciones de roles para una entidad específica
func (uc *UseCase) GetByEntity(ctx context.Context, entityType string, entityID int64) ([]*AssignedRole, error) {
	return uc.repo.GetByEntity(ctx, entityType, entityID)
}

// Assign asigna un rol a una entidad
func (uc *UseCase) Assign(ctx context.Context, assignedRole *AssignedRole) error {
	// Verificar que el rol exista
	_, err := uc.roleRepo.GetByID(ctx, assignedRole.RoleID)
	if err != nil {
		return err
	}

	return uc.repo.Assign(ctx, assignedRole)
}

// Revoke revoca un rol de una entidad
func (uc *UseCase) Revoke(ctx context.Context, roleID, entityID int64, entityType string) error {
	return uc.repo.Revoke(ctx, roleID, entityID, entityType)
}

// HasRole verifica si una entidad tiene un rol específico
func (uc *UseCase) HasRole(ctx context.Context, roleID, entityID int64, entityType string) (bool, error) {
	return uc.repo.HasRole(ctx, roleID, entityID, entityType)
}

// List obtiene una lista paginada de asignaciones de roles con filtros opcionales
func (uc *UseCase) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*AssignedRole, int, error) {
	return uc.repo.List(ctx, filters, page, pageSize)
}
