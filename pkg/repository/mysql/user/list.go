package user

import (
	"context"
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/user"
)

// List lista usuarios con paginación y filtros
func (r *Repository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*user.User, int, error) {
	// Construir la consulta base
	baseQuery := `
		SELECT id, name, email, password, role_id,
		       email_verified_at, remember_token, created_at, updated_at, deleted_at
		FROM users
		WHERE deleted_at IS NULL
	`

	countQuery := `
		SELECT COUNT(*)
		FROM users
		WHERE deleted_at IS NULL
	`

	// Aplicar filtros si existen
	args := []interface{}{}
	whereClause := ""

	if name, ok := filters["name"].(string); ok && name != "" {
		whereClause += " AND name LIKE ?"
		args = append(args, "%"+name+"%")
	}

	if email, ok := filters["email"].(string); ok && email != "" {
		whereClause += " AND email LIKE ?"
		args = append(args, "%"+email+"%")
	}

	if roleID, ok := filters["role_id"].(string); ok && roleID != "" {
		whereClause += " AND role_id = ?"
		args = append(args, roleID)
	}

	// Aplicar paginación
	offset := (page - 1) * pageSize
	query := baseQuery + whereClause + " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	countQueryFinal := countQuery + whereClause

	// Argumentos para la paginación
	queryArgs := append(args, pageSize, offset)

	// Obtener el total de registros
	var total int
	err := r.db.QueryRowContext(ctx, countQueryFinal, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Ejecutar la consulta paginada
	rows, err := r.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	// Procesar resultados
	users := []*user.User{}
	for rows.Next() {
		var u user.User
		var emailVerifiedAt, createdAt, updatedAt, deletedAt sql.NullTime
		var rememberToken sql.NullString

		err := rows.Scan(
			&u.ID, &u.Name, &u.Email, &u.Password, &u.RoleID,
			&emailVerifiedAt, &rememberToken, &createdAt, &updatedAt, &deletedAt,
		)

		if err != nil {
			return nil, 0, err
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

		users = append(users, &u)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
