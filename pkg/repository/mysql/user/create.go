package user

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/your-org/jvairv2/pkg/domain/user"
)

// Create crea un nuevo usuario
func (r *Repository) Create(ctx context.Context, u *user.User) error {
	slog.Debug("Creando usuario en repositorio",
		"name", u.Name,
		"email", u.Email,
	)

	// Hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Error al generar hash de contraseña",
			"error", err,
		)
		return fmt.Errorf("error al generar hash de contraseña: %w", err)
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

	// Manejar el role_id que puede ser nulo
	var roleIDValue interface{}
	if u.RoleID != nil {
		roleIDValue = *u.RoleID
	} else {
		roleIDValue = nil
	}

	now := time.Now()
	nowPtr := &now

	result, err := r.db.ExecContext(ctx, query,
		u.Name, u.Email, string(hashedPassword), roleIDValue,
		emailVerifiedAtValue, nowPtr, nowPtr,
	)

	if err != nil {
		slog.Error("Error en la inserción de usuario",
			"email", u.Email,
			"error", err,
		)
		return fmt.Errorf("error en la inserción de usuario: %w", err)
	}

	// Obtener el ID generado
	lastID, err := result.LastInsertId()
	if err != nil {
		slog.Error("Error al obtener el ID generado",
			"error", err,
		)
		return fmt.Errorf("error al obtener el ID generado: %w", err)
	}

	u.ID = lastID
	u.CreatedAt = nowPtr
	u.UpdatedAt = nowPtr

	slog.Debug("Usuario creado exitosamente en repositorio",
		"user_id", u.ID,
		"email", u.Email,
	)
	return nil
}
