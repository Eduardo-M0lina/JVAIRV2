package user

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/user"
)

// GetByEmail obtiene un usuario por su email
func (r *Repository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	slog.Debug("Buscando usuario por email",
		"email", email,
	)

	if r.db == nil {
		slog.Error("Conexión a base de datos no inicializada")
		return nil, fmt.Errorf("conexión a base de datos no inicializada")
	}

	query := `
		SELECT u.id, u.name, u.email, u.password, u.role_id,
		       u.email_verified_at, u.remember_token, u.created_at, u.updated_at, u.deleted_at,
		       r.id as role_id_int, r.name as role_name, r.title as role_title
		FROM users u
		LEFT JOIN (
			SELECT ar.entity_id, ar.role_id
			FROM assigned_roles ar
			WHERE ar.entity_type = 'App\\Models\\User'
			GROUP BY ar.entity_id
		) ar ON ar.entity_id = u.id
		LEFT JOIN roles r ON r.id = ar.role_id
		WHERE u.email = ? AND u.deleted_at IS NULL
	`

	var u user.User
	var emailVerifiedAt, createdAt, updatedAt, deletedAt sql.NullTime
	var rememberToken sql.NullString
	var roleID sql.NullString
	var roleIDInt sql.NullInt64
	var roleName, roleTitle sql.NullString

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.Name, &u.Email, &u.Password, &roleID,
		&emailVerifiedAt, &rememberToken, &createdAt, &updatedAt, &deletedAt,
		&roleIDInt, &roleName, &roleTitle,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			slog.Debug("Usuario no encontrado con email",
				"email", email,
			)
			return nil, ErrUserNotFound
		}
		slog.Error("Error al consultar usuario por email",
			"email", email,
			"error", err,
		)
		return nil, fmt.Errorf("error al consultar usuario por email: %w", err)
	}

	slog.Debug("Usuario encontrado por email",
		"email", email,
		"user_id", u.ID,
	)

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

	// Información del rol
	if roleName.Valid {
		u.RoleName = &roleName.String
	}
	if roleTitle.Valid {
		u.RoleTitle = &roleTitle.String
	}

	// Campo virtual
	u.IsActive = u.DeletedAt == nil

	return &u, nil
}
