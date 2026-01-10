package role

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/your-org/jvairv2/pkg/domain/role"
)

// List obtiene una lista paginada de roles con filtros opcionales
func (r *Repository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*role.Role, int, error) {
	// Construir la consulta base
	countQuery := "SELECT COUNT(*) FROM roles"
	selectQuery := `
		SELECT id, name, title, scope, created_at, updated_at
		FROM roles
	`

	// Aplicar filtros si existen
	whereConditions := []string{}
	args := []interface{}{}

	if name, ok := filters["name"].(string); ok && name != "" {
		whereConditions = append(whereConditions, "name LIKE ?")
		args = append(args, "%"+name+"%")
	}

	if title, ok := filters["title"].(string); ok && title != "" {
		whereConditions = append(whereConditions, "title LIKE ?")
		args = append(args, "%"+title+"%")
	}

	if scope, ok := filters["scope"].(int); ok {
		whereConditions = append(whereConditions, "scope = ?")
		args = append(args, scope)
	}

	// Construir la cláusula WHERE si hay condiciones
	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = " WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Aplicar la cláusula WHERE a las consultas
	countQueryWithFilters := countQuery + whereClause
	selectQueryWithFilters := selectQuery + whereClause

	// Obtener el total de registros
	var total int
	err := r.db.QueryRowContext(ctx, countQueryWithFilters, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Si no hay resultados, devolver una lista vacía
	if total == 0 {
		return []*role.Role{}, 0, nil
	}

	// Aplicar paginación
	offset := (page - 1) * pageSize
	selectQueryWithPagination := fmt.Sprintf("%s ORDER BY name LIMIT %d OFFSET %d", selectQueryWithFilters, pageSize, offset)

	// Ejecutar la consulta paginada
	rows, err := r.db.QueryContext(ctx, selectQueryWithPagination, args...)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	// Procesar resultados
	roles := []*role.Role{}
	for rows.Next() {
		var r role.Role
		var title sql.NullString
		var scope sql.NullInt32
		var createdAt, updatedAt sql.NullTime

		err := rows.Scan(&r.ID, &r.Name, &title, &scope, &createdAt, &updatedAt)
		if err != nil {
			return nil, 0, err
		}

		if title.Valid {
			r.Title = &title.String
		}

		if scope.Valid {
			scopeInt := int(scope.Int32)
			r.Scope = &scopeInt
		}

		if createdAt.Valid {
			r.CreatedAt = &createdAt.Time
		}

		if updatedAt.Valid {
			r.UpdatedAt = &updatedAt.Time
		}

		roles = append(roles, &r)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}
