package ability

import (
	"context"
)

// Repository define las operaciones de persistencia para abilities (capacidades/permisos)
type Repository interface {
	// Obtener una ability por ID
	GetByID(ctx context.Context, id int64) (*Ability, error)

	// Obtener una ability por nombre
	GetByName(ctx context.Context, name string) (*Ability, error)

	// Crear una nueva ability
	Create(ctx context.Context, ability *Ability) error

	// Actualizar una ability existente
	Update(ctx context.Context, ability *Ability) error

	// Eliminar una ability
	Delete(ctx context.Context, id int64) error

	// Listar abilities con paginación y filtros
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Ability, int, error)
}
