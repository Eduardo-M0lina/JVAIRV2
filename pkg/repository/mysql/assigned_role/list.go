package assigned_role

import (
	"context"
	"database/sql"
	"strings"

	domainAssignedRole "github.com/your-org/jvairv2/pkg/domain/assigned_role"
)

// List obtiene una lista paginada de asignaciones de roles con filtros opcionales
func (r *Repository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*domainAssignedRole.AssignedRole, int, error) {
	// Construir la consulta base
	countQuery := "SELECT COUNT(*) FROM assigned_roles"
	selectQuery := `
		SELECT id, role_id, entity_id, entity_type, restricted_to_id, restricted_to_type, scope
		FROM assigned_roles
	`

	// Procesar filtros
	var whereConditions []string
	var args []interface{}

	for key, value := range filters {
		switch key {
		case "role_id":
			whereConditions = append(whereConditions, "role_id = ?")
			args = append(args, value)
		case "entity_type":
			whereConditions = append(whereConditions, "entity_type = ?")
			args = append(args, value)
		case "entity_id":
			whereConditions = append(whereConditions, "entity_id = ?")
			args = append(args, value)
		case "restricted":
			whereConditions = append(whereConditions, "restricted = ?")
			args = append(args, value)
		}
	}

	// Agregar condiciones WHERE si existen
	if len(whereConditions) > 0 {
		whereClause := " WHERE " + strings.Join(whereConditions, " AND ")
		countQuery += whereClause
		selectQuery += whereClause
	}

	// Agregar paginación
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	selectQuery += " ORDER BY id ASC LIMIT ? OFFSET ?"

	// Ejecutar consulta de conteo
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Si no hay resultados, devolver un slice vacío
	if total == 0 {
		return []*domainAssignedRole.AssignedRole{}, 0, nil
	}

	// Agregar argumentos de paginación
	args = append(args, pageSize, offset)

	// Ejecutar consulta de selección
	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	// Procesar resultados
	var assignedRoles []*domainAssignedRole.AssignedRole
	for rows.Next() {
		var assignedRole domainAssignedRole.AssignedRole
		var restrictedToID sql.NullInt64
		var restrictedToType sql.NullString
		var scope sql.NullInt32

		err := rows.Scan(
			&assignedRole.ID, &assignedRole.RoleID, &assignedRole.EntityID, &assignedRole.EntityType,
			&restrictedToID, &restrictedToType, &scope,
		)
		if err != nil {
			return nil, 0, err
		}

		if restrictedToID.Valid {
			assignedRole.RestrictedToID = &restrictedToID.Int64
			assignedRole.Restricted = true
		}

		if restrictedToType.Valid {
			assignedRole.RestrictedToType = &restrictedToType.String
		}

		if scope.Valid {
			scopeInt := int(scope.Int32)
			assignedRole.Scope = &scopeInt
		}

		assignedRoles = append(assignedRoles, &assignedRole)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return assignedRoles, total, nil
}
