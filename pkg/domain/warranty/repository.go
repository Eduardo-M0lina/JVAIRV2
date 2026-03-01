package warranty

import "context"

// Repository define los métodos para interactuar con el almacenamiento de warranties
type Repository interface {
	// Create crea una nueva garantía
	Create(ctx context.Context, w *Warranty) error

	// GetByID obtiene una garantía por su ID
	GetByID(ctx context.Context, id int64) (*Warranty, error)

	// List obtiene una lista paginada de garantías con filtros opcionales
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Warranty, int, error)

	// Update actualiza una garantía existente
	Update(ctx context.Context, w *Warranty) error

	// Delete elimina una garantía (soft delete)
	Delete(ctx context.Context, id int64) error
}
