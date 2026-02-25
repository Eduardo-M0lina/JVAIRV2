package invoice

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	domainInvoice "github.com/your-org/jvairv2/pkg/domain/invoice"
)

// List obtiene una lista paginada de facturas con filtros opcionales
func (r *Repository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*domainInvoice.Invoice, int, error) {
	var conditions []string
	var args []interface{}
	var havingConditions []string

	// Siempre excluir soft-deleted
	conditions = append(conditions, "i.deleted_at IS NULL")

	// Filtro por job_id
	if jobID, ok := filters["job_id"].(int64); ok && jobID > 0 {
		conditions = append(conditions, "i.job_id = ?")
		args = append(args, jobID)
	}

	// Búsqueda en múltiples campos (fiel al original Laravel: invoice_number, work_order, property, customer)
	if search, ok := filters["search"].(string); ok && search != "" {
		searchCondition := `(
			i.invoice_number LIKE ? OR
			j.work_order LIKE ? OR
			p.property_code LIKE ? OR
			p.street LIKE ? OR
			p.city LIKE ? OR
			p.state LIKE ? OR
			p.zip LIKE ? OR
			c.name LIKE ? OR
			c.email LIKE ? OR
			c.phone LIKE ? OR
			c.mobile LIKE ?
		)`
		conditions = append(conditions, searchCondition)
		searchPattern := "%" + search + "%"
		for i := 0; i < 11; i++ {
			args = append(args, searchPattern)
		}
	}

	// Filtro por status (paid/unpaid) — se aplica como HAVING sobre balance calculado
	if status, ok := filters["status"].(string); ok && status != "" {
		switch status {
		case "paid":
			havingConditions = append(havingConditions, "balance <= 0")
		case "unpaid":
			havingConditions = append(havingConditions, "balance > 0")
		}
	}

	whereClause := strings.Join(conditions, " AND ")

	havingClause := ""
	if len(havingConditions) > 0 {
		havingClause = "HAVING " + strings.Join(havingConditions, " AND ")
	}

	// Count query — necesita subquery para contar correctamente con GROUP BY + HAVING
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM (
			SELECT i.id,
				IFNULL(i.total - SUM(ip.amount), i.total) as balance
			FROM invoices i
			LEFT JOIN invoice_payments ip ON ip.invoice_id = i.id AND ip.deleted_at IS NULL
			LEFT JOIN jobs j ON j.id = i.job_id
			LEFT JOIN properties p ON p.id = j.property_id
			LEFT JOIN customers c ON c.id = p.customer_id
			WHERE %s
			GROUP BY i.id
			%s
		) sub
	`, whereClause, havingClause)

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		slog.ErrorContext(ctx, "Failed to count invoices",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	// Sorting
	orderClause := "i.created_at DESC"
	if sort, ok := filters["sort"].(string); ok && sort != "" {
		direction := "DESC"
		if dir, ok := filters["direction"].(string); ok && strings.ToUpper(dir) == "ASC" {
			direction = "ASC"
		}

		switch sort {
		case "invoice_number":
			orderClause = fmt.Sprintf("i.invoice_number %s", direction)
		case "total":
			orderClause = fmt.Sprintf("i.total %s", direction)
		case "balance":
			orderClause = fmt.Sprintf("balance %s", direction)
		case "created_at":
			orderClause = fmt.Sprintf("i.created_at %s", direction)
		}
	}

	// Data query
	offset := (page - 1) * pageSize
	dataQuery := fmt.Sprintf(`
		SELECT
			i.id, i.job_id, i.invoice_number, i.total, i.description,
			i.allow_online_payments, i.notes,
			i.created_at, i.updated_at, i.deleted_at,
			IFNULL(i.total - SUM(ip.amount), i.total) as balance
		FROM invoices i
		LEFT JOIN invoice_payments ip ON ip.invoice_id = i.id AND ip.deleted_at IS NULL
		LEFT JOIN jobs j ON j.id = i.job_id
		LEFT JOIN properties p ON p.id = j.property_id
		LEFT JOIN customers c ON c.id = p.customer_id
		WHERE %s
		GROUP BY i.id
		%s
		ORDER BY %s
		LIMIT ? OFFSET ?
	`, whereClause, havingClause, orderClause)

	queryArgs := append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, queryArgs...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list invoices",
			slog.String("error", err.Error()))
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var invoices []*domainInvoice.Invoice
	for rows.Next() {
		inv := &domainInvoice.Invoice{}
		var balance float64
		if err := rows.Scan(
			&inv.ID, &inv.JobID, &inv.InvoiceNumber, &inv.Total, &inv.Description,
			&inv.AllowOnlinePayments, &inv.Notes,
			&inv.CreatedAt, &inv.UpdatedAt, &inv.DeletedAt,
			&balance,
		); err != nil {
			slog.ErrorContext(ctx, "Failed to scan invoice row",
				slog.String("error", err.Error()))
			return nil, 0, err
		}
		inv.Balance = &balance
		invoices = append(invoices, inv)
	}

	if err = rows.Err(); err != nil {
		slog.ErrorContext(ctx, "Error iterating invoice rows",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	return invoices, total, nil
}
