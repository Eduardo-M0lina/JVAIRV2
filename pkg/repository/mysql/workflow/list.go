package workflow

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	domainWorkflow "github.com/your-org/jvairv2/pkg/domain/workflow"
)

// List obtiene una lista paginada de workflows
func (r *Repository) List(ctx context.Context, filters domainWorkflow.Filters, page, pageSize int) ([]domainWorkflow.Workflow, int64, error) {
	// Construir la consulta base
	var conditions []string
	var args []interface{}

	// Aplicar filtros
	if filters.Name != "" {
		conditions = append(conditions, "name LIKE ?")
		args = append(args, "%"+filters.Name+"%")
	}

	if filters.IsActive != nil {
		conditions = append(conditions, "is_active = ?")
		args = append(args, *filters.IsActive)
	}

	if filters.Search != "" {
		conditions = append(conditions, "(name LIKE ? OR notes LIKE ?)")
		searchTerm := "%" + filters.Search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	// Construir la cláusula WHERE
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Contar el total de registros
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM workflows %s", whereClause)
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Calcular el offset
	offset := (page - 1) * pageSize

	// Consulta para obtener los workflows
	selectQuery := fmt.Sprintf(`
		SELECT id, name, notes, is_active, created_at, updated_at
		FROM workflows
		%s
		ORDER BY id DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var workflows []domainWorkflow.Workflow
	for rows.Next() {
		var workflow domainWorkflow.Workflow
		var notes sql.NullString
		var createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&workflow.ID,
			&workflow.Name,
			&notes,
			&workflow.IsActive,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		if notes.Valid {
			workflow.Notes = &notes.String
		}
		if createdAt.Valid {
			workflow.CreatedAt = &createdAt.Time
		}
		if updatedAt.Valid {
			workflow.UpdatedAt = &updatedAt.Time
		}

		workflows = append(workflows, workflow)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return workflows, total, nil
}
