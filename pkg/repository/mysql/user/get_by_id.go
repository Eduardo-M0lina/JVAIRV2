package user

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/your-org/jvairv2/pkg/domain/user"
)

// GetByID obtiene un usuario por su ID
func (r *Repository) GetByID(ctx context.Context, id string) (*user.User, error) {
	// Convertir ID de string a int64
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, errors.New("ID de usuario inválido")
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
		WHERE u.id = ? AND u.deleted_at IS NULL
	`

	var u user.User
	var emailVerifiedAt, createdAt, updatedAt, deletedAt sql.NullTime
	var rememberToken sql.NullString
	var roleIDInt sql.NullInt64
	var roleName, roleTitle sql.NullString

	err = r.db.QueryRowContext(ctx, query, idInt).Scan(
		&u.ID, &u.Name, &u.Email, &u.Password, &u.RoleID,
		&emailVerifiedAt, &rememberToken, &createdAt, &updatedAt, &deletedAt,
		&roleIDInt, &roleName, &roleTitle,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
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
