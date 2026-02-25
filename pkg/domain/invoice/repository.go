package invoice

import "context"

// Repository define los métodos para interactuar con el almacenamiento de invoices
type Repository interface {
	// Create crea una nueva factura
	Create(ctx context.Context, inv *Invoice) error

	// GetByID obtiene una factura por su ID (incluye balance calculado)
	GetByID(ctx context.Context, id int64) (*Invoice, error)

	// List obtiene una lista paginada de facturas con filtros opcionales
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Invoice, int, error)

	// Update actualiza una factura existente
	Update(ctx context.Context, inv *Invoice) error

	// Delete elimina una factura (soft delete)
	Delete(ctx context.Context, id int64) error
}
