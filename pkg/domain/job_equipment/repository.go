package job_equipment

import "context"

// Repository define los métodos para interactuar con el almacenamiento de equipos de trabajo
type Repository interface {
	// Create crea un nuevo equipo de trabajo
	Create(ctx context.Context, equipment *JobEquipment) error

	// GetByID obtiene un equipo de trabajo por su ID
	GetByID(ctx context.Context, id int64) (*JobEquipment, error)

	// List obtiene una lista de equipos de un trabajo con filtro opcional por type
	List(ctx context.Context, jobID int64, equipmentType string) ([]*JobEquipment, error)

	// Update actualiza un equipo de trabajo existente
	Update(ctx context.Context, equipment *JobEquipment) error

	// Delete elimina un equipo de trabajo (hard delete)
	Delete(ctx context.Context, id int64) error
}
