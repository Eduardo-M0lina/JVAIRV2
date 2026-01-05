package role

import (
	"context"
)

// Repository define las operaciones de persistencia para roles
type Repository interface {
	// Obtener un rol por ID
	GetByID(ctx context.Context, id int64) (*Role, error)

	// Obtener un rol por nombre
	GetByName(ctx context.Context, name string) (*Role, error)

	// Crear un nuevo rol
	Create(ctx context.Context, role *Role) error

	// Actualizar un rol existente
	Update(ctx context.Context, role *Role) error

	// Eliminar un rol
	Delete(ctx context.Context, id int64) error

	// Listar roles con paginación y filtros
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Role, int, error)
}
