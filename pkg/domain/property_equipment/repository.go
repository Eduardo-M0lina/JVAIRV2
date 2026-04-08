package property_equipment

import "context"

// Repository define los métodos para interactuar con el almacenamiento de equipos de propiedad
type Repository interface {
	// Create crea un nuevo equipo de propiedad
	Create(ctx context.Context, equipment *PropertyEquipment) error

	// GetByID obtiene un equipo de propiedad por su ID
	GetByID(ctx context.Context, id int64) (*PropertyEquipment, error)

	// List obtiene una lista de equipos de una propiedad
	List(ctx context.Context, propertyID int64) ([]*PropertyEquipment, error)

	// Update actualiza un equipo de propiedad existente
	Update(ctx context.Context, equipment *PropertyEquipment) error

	// Delete elimina un equipo de propiedad (hard delete)
	Delete(ctx context.Context, id int64) error
}
