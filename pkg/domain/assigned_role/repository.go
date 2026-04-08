package assigned_role

import (
	"context"
)

// Repository define las operaciones de persistencia para asignaciones de roles
type Repository interface {
	// Obtener una asignación de rol por ID
	GetByID(ctx context.Context, id int64) (*AssignedRole, error)

	// Obtener asignaciones de rol por entidad
	GetByEntity(ctx context.Context, entityType string, entityID int64) ([]*AssignedRole, error)

	// Asignar un rol a una entidad
	Assign(ctx context.Context, assignedRole *AssignedRole) error

	// Revocar un rol de una entidad
	Revoke(ctx context.Context, roleID, entityID int64, entityType string) error

	// Verificar si una entidad tiene un rol específico
	HasRole(ctx context.Context, roleID, entityID int64, entityType string) (bool, error)

	// Listar asignaciones de roles con paginación y filtros
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*AssignedRole, int, error)
}
