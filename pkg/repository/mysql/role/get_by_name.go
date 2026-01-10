package role

import (
	"context"
	"database/sql"

	domainRole "github.com/your-org/jvairv2/pkg/domain/role"
)

// GetByName obtiene un rol por su nombre
func (r *Repository) GetByName(ctx context.Context, name string) (*domainRole.Role, error) {
	query := `
		SELECT id, name, title, scope, created_at, updated_at
		FROM roles
		WHERE name = ?
	`

	var roleEntity domainRole.Role
	var title sql.NullString
	var scope sql.NullInt32
	var createdAt, updatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&roleEntity.ID, &roleEntity.Name, &title, &scope, &createdAt, &updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}

	if title.Valid {
		roleEntity.Title = &title.String
	}

	if scope.Valid {
		scopeInt := int(scope.Int32)
		roleEntity.Scope = &scopeInt
	}

	if createdAt.Valid {
		roleEntity.CreatedAt = &createdAt.Time
	}

	if updatedAt.Valid {
		roleEntity.UpdatedAt = &updatedAt.Time
	}

	return &roleEntity, nil
}
