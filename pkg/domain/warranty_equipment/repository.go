package warranty_equipment

import "context"

// Repository define los métodos para interactuar con el almacenamiento de equipos de garantía
type Repository interface {
	// Create crea un nuevo equipo de garantía
	Create(ctx context.Context, equipment *WarrantyEquipment) error

	// GetByID obtiene un equipo de garantía por su ID
	GetByID(ctx context.Context, id int64) (*WarrantyEquipment, error)

	// ListByWarrantyID obtiene una lista de equipos de una garantía
	ListByWarrantyID(ctx context.Context, warrantyID int64) ([]*WarrantyEquipment, error)

	// Update actualiza un equipo de garantía existente
	Update(ctx context.Context, equipment *WarrantyEquipment) error

	// Delete elimina un equipo de garantía (hard delete)
	Delete(ctx context.Context, id int64) error

	// CloneFromJobEquipment clona los equipos de un job a una garantía
	CloneFromJobEquipment(ctx context.Context, warrantyID int64, jobID int64) error
}
