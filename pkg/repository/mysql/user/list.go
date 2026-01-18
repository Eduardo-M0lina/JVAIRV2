package user

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/user"
)

// List lista usuarios con paginación y filtros
func (r *Repository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*user.User, int, error) {
	// Construir la consulta base con JOIN a assigned_roles y roles
	// Usamos una subconsulta para obtener solo el primer rol asignado por usuario
	baseQuery := `
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
		WHERE u.deleted_at IS NULL
	`

	countQuery := `
		SELECT COUNT(*)
		FROM users u
		WHERE u.deleted_at IS NULL
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
	slog.Debug("Ejecutando query de listado de usuarios",
		"query", query,
		"args", queryArgs,
	)
	rows, err := r.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		slog.Error("Error al ejecutar query de listado",
			"error", err,
		)
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	// Procesar resultados
	users := []*user.User{}
	for rows.Next() {
		var u user.User
		var emailVerifiedAt, createdAt, updatedAt, deletedAt sql.NullTime
		var rememberToken sql.NullString
		var roleIDInt sql.NullInt64
		var roleName, roleTitle sql.NullString

		err := rows.Scan(
			&u.ID, &u.Name, &u.Email, &u.Password, &u.RoleID,
			&emailVerifiedAt, &rememberToken, &createdAt, &updatedAt, &deletedAt,
			&roleIDInt, &roleName, &roleTitle,
		)

		if err != nil {
			slog.Error("Error al escanear fila de usuario",
				"error", err,
			)
			return nil, 0, err
		}

		slog.Debug("Usuario escaneado",
			"user_id", u.ID,
			"email", u.Email,
			"role_id_field", u.RoleID,
			"role_id_int", roleIDInt,
			"role_name", roleName,
			"role_title", roleTitle,
		)

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
			slog.Debug("Rol asignado al usuario",
				"user_id", u.ID,
				"role_name", *u.RoleName,
			)
		} else {
			slog.Debug("Usuario sin rol asignado",
				"user_id", u.ID,
				"email", u.Email,
			)
		}
		if roleTitle.Valid {
			u.RoleTitle = &roleTitle.String
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
