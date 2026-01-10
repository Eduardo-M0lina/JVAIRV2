package user

import (
	"context"
	"strconv"
	"time"

	"github.com/your-org/jvairv2/pkg/domain/user"
)

// Update actualiza un usuario existente
func (r *Repository) Update(ctx context.Context, u *user.User) error {
	// Verificar si el usuario existe
	idStr := strconv.FormatInt(u.ID, 10)
	_, err := r.GetByID(ctx, idStr)
	if err != nil {
		return err
	}

	query := `
		UPDATE users
		SET name = ?, email = ?, role_id = ?,
		    email_verified_at = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	now := time.Now()
	u.UpdatedAt = &now

	var emailVerifiedAtValue interface{}
	if u.EmailVerifiedAt != nil {
		emailVerifiedAtValue = *u.EmailVerifiedAt
	} else {
		emailVerifiedAtValue = nil
	}

	_, err = r.db.ExecContext(ctx, query,
		u.Name, u.Email, u.RoleID,
		emailVerifiedAtValue, now, u.ID,
	)

	return err
}
