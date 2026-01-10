package user

import (
	"context"

	"github.com/your-org/jvairv2/pkg/domain/user"
	"golang.org/x/crypto/bcrypt"
)

// VerifyCredentials verifica las credenciales de un usuario
func (r *Repository) VerifyCredentials(ctx context.Context, email, password string) (*user.User, error) {
	u, err := r.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Verificar si el usuario está activo
	if !u.IsActive {
		return nil, ErrInvalidCredentials
	}

	// Verificar contraseña
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return u, nil
}
