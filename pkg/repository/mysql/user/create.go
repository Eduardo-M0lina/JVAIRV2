package user

import (
	"context"
	"time"

	"github.com/your-org/jvairv2/pkg/domain/user"
	"golang.org/x/crypto/bcrypt"
)

// Create crea un nuevo usuario
func (r *Repository) Create(ctx context.Context, u *user.User) error {
	// Verificar si el email ya existe
	_, err := r.GetByEmail(ctx, u.Email)
	if err == nil {
		return ErrDuplicateEmail
	} else if err != ErrUserNotFound {
		return err
	}

	// Hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO users (name, email, password, role_id,
		                  email_verified_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	var emailVerifiedAtValue interface{}
	if u.EmailVerifiedAt != nil {
		emailVerifiedAtValue = *u.EmailVerifiedAt
	} else {
		emailVerifiedAtValue = nil
	}

	now := time.Now()
	nowPtr := &now

	result, err := r.db.ExecContext(ctx, query,
		u.Name, u.Email, string(hashedPassword), u.RoleID,
		emailVerifiedAtValue, nowPtr, nowPtr,
	)

	if err != nil {
		return err
	}

	// Obtener el ID generado
	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = lastID
	u.CreatedAt = nowPtr
	u.UpdatedAt = nowPtr

	return nil
}
