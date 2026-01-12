package settings

import (
	"context"
	"errors"
)

var (
	// ErrSettingsNotFound se devuelve cuando no se encuentran las configuraciones
	ErrSettingsNotFound = errors.New("configuraciones no encontradas")
)

// UseCase maneja la lógica de negocio de las configuraciones
type UseCase struct {
	repo Repository
}

// NewUseCase crea una nueva instancia del caso de uso de configuraciones
func NewUseCase(repo Repository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

// Get obtiene las configuraciones del sistema
func (uc *UseCase) Get(ctx context.Context) (*Settings, error) {
	return uc.repo.Get(ctx)
}

// Update actualiza las configuraciones del sistema
func (uc *UseCase) Update(ctx context.Context, settings *Settings) error {
	// Validar que los valores sean correctos
	if settings.PasswordExpireDays < 1 {
		return errors.New("los días de expiración de contraseña deben ser al menos 1")
	}

	if settings.PasswordHistoryCount < 0 {
		return errors.New("el conteo de historial de contraseñas no puede ser negativo")
	}

	if settings.PasswordMinimumLength < 4 {
		return errors.New("la longitud mínima de contraseña debe ser al menos 4")
	}

	if settings.PasswordAge < 0 {
		return errors.New("la edad de contraseña no puede ser negativa")
	}

	return uc.repo.Update(ctx, settings)
}
