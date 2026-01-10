package permission

import (
	"context"
)

// Repository define las operaciones de persistencia para permisos
type Repository interface {
	// Obtener un permiso por ID
	GetByID(ctx context.Context, id int64) (*Permission, error)

	// Obtener permisos por entidad
	GetByEntity(ctx context.Context, entityType string, entityID int64) ([]*Permission, error)

	// Obtener permisos por ability
	GetByAbility(ctx context.Context, abilityID int64) ([]*Permission, error)

	// Crear un nuevo permiso
	Create(ctx context.Context, permission *Permission) error

	// Actualizar un permiso existente
	Update(ctx context.Context, permission *Permission) error

	// Eliminar un permiso
	Delete(ctx context.Context, id int64) error

	// Verificar si existe un permiso específico
	Exists(ctx context.Context, abilityID, entityID int64, entityType string) (bool, error)

	// Listar permisos con paginación y filtros
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Permission, int, error)
}
