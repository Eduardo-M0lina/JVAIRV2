package user

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/your-org/jvairv2/pkg/domain/user"
)

// GetByEmail obtiene un usuario por su email
func (r *Repository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	log.Printf("[REPO] Buscando usuario por email: %s", email)

	if r.db == nil {
		log.Printf("[REPO] ERROR: Conexión a base de datos no inicializada")
		return nil, fmt.Errorf("conexión a base de datos no inicializada")
	}

	query := `
		SELECT id, name, email, password, role_id,
		       email_verified_at, remember_token, created_at, updated_at, deleted_at
		FROM users
		WHERE email = ? AND deleted_at IS NULL
	`

	var u user.User
	var emailVerifiedAt, createdAt, updatedAt, deletedAt sql.NullTime
	var rememberToken sql.NullString
	var roleID sql.NullString

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.Name, &u.Email, &u.Password, &roleID,
		&emailVerifiedAt, &rememberToken, &createdAt, &updatedAt, &deletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[REPO] Usuario no encontrado con email: %s", email)
			return nil, ErrUserNotFound
		}
		log.Printf("[REPO] ERROR al consultar usuario por email %s: %v", email, err)
		return nil, fmt.Errorf("error al consultar usuario por email: %w", err)
	}

	log.Printf("[REPO] Usuario encontrado con email %s, ID: %d", email, u.ID)

	// Asignar roleID si es válido
	if roleID.Valid {
		u.RoleID = &roleID.String
	}

	if emailVerifiedAt.Valid {
		u.EmailVerifiedAt = &emailVerifiedAt.Time
	}

	if rememberToken.Valid {
		u.RememberToken = &rememberToken.String
	}

	if createdAt.Valid {
		u.CreatedAt = &createdAt.Time
	}

	if updatedAt.Valid {
		u.UpdatedAt = &updatedAt.Time
	}

	if deletedAt.Valid {
		u.DeletedAt = &deletedAt.Time
	}

	// Campo virtual
	u.IsActive = u.DeletedAt == nil

	log.Printf("[REPO] Retornando usuario con email %s, ID: %d", u.Email, u.ID)
	return &u, nil
}
