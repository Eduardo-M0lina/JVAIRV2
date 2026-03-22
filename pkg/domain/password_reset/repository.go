package password_reset

import "context"

// Repository define las operaciones de persistencia para tokens de reseteo de contraseña
type Repository interface {
	// Create crea un nuevo token de reseteo
	Create(ctx context.Context, pr *PasswordReset) error

	// GetByToken obtiene un token de reseteo por su valor
	GetByToken(ctx context.Context, token string) (*PasswordReset, error)

	// GetByEmail obtiene el token más reciente por email
	GetByEmail(ctx context.Context, email string) (*PasswordReset, error)

	// Delete elimina un token de reseteo
	Delete(ctx context.Context, token string) error

	// DeleteByEmail elimina todos los tokens de un email
	DeleteByEmail(ctx context.Context, email string) error

	// DeleteExpired elimina tokens expirados (más antiguos que X horas)
	DeleteExpired(ctx context.Context, hours int) error
}
