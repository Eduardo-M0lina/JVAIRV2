package password_history

import "context"

// Repository define las operaciones de persistencia para historial de contraseñas
type Repository interface {
	// Create crea un nuevo registro de historial
	Create(ctx context.Context, ph *PasswordHistory) error

	// GetByUserID obtiene el historial de contraseñas de un usuario
	GetByUserID(ctx context.Context, userID int64, limit int) ([]*PasswordHistory, error)

	// Delete elimina registros de historial por ID
	Delete(ctx context.Context, id int64) error

	// DeleteOldest elimina los registros más antiguos de un usuario, manteniendo solo los N más recientes
	DeleteOldest(ctx context.Context, userID int64, keep int) error
}
