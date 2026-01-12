package settings

import "context"

// Repository define los métodos para interactuar con la base de datos de configuraciones
type Repository interface {
	Get(ctx context.Context) (*Settings, error)
	Update(ctx context.Context, settings *Settings) error
}
