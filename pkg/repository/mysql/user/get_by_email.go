package user

import (
	"context"
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/user"
)

// GetByEmail obtiene un usuario por su email
func (r *Repository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT id, name, email, password, role_id,
		       email_verified_at, remember_token, created_at, updated_at, deleted_at
		FROM users
		WHERE email = ? AND deleted_at IS NULL
	`

	var u user.User
	var emailVerifiedAt, createdAt, updatedAt, deletedAt sql.NullTime
	var rememberToken sql.NullString

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.Name, &u.Email, &u.Password, &u.RoleID,
		&emailVerifiedAt, &rememberToken, &createdAt, &updatedAt, &deletedAt,
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

	// Campo virtual
	u.IsActive = u.DeletedAt == nil

	return &u, nil
}
