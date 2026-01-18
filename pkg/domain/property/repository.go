package property

import "context"

// Repository define los métodos para interactuar con el almacenamiento de propiedades
type Repository interface {
	// Create crea una nueva propiedad
	Create(ctx context.Context, property *Property) error

	// GetByID obtiene una propiedad por su ID
	GetByID(ctx context.Context, id int64) (*Property, error)

	// List obtiene una lista paginada de propiedades con filtros opcionales
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Property, int, error)

	// Update actualiza una propiedad existente
	Update(ctx context.Context, property *Property) error

	// Delete elimina una propiedad (soft delete)
	Delete(ctx context.Context, id int64) error

	// HasJobs verifica si una propiedad tiene jobs asociados
	HasJobs(ctx context.Context, id int64) (bool, error)
}
