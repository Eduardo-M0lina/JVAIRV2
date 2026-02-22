package job

import "context"

// Repository define los métodos para interactuar con el almacenamiento de jobs
type Repository interface {
	// Create crea un nuevo job
	Create(ctx context.Context, job *Job) error

	// GetByID obtiene un job por su ID
	GetByID(ctx context.Context, id int64) (*Job, error)

	// List obtiene una lista paginada de jobs con filtros opcionales
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Job, int, error)

	// Update actualiza un job existente
	Update(ctx context.Context, job *Job) error

	// Delete elimina un job (soft delete)
	Delete(ctx context.Context, id int64) error

	// Close cierra un job
	Close(ctx context.Context, id int64, jobStatusID int64) error
}
