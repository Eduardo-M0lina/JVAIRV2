package invoice_payment

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	domainPayment "github.com/your-org/jvairv2/pkg/domain/invoice_payment"
)

// ListByInvoiceID obtiene los pagos de una factura con paginación
func (r *Repository) ListByInvoiceID(ctx context.Context, invoiceID int64, filters map[string]interface{}, page, pageSize int) ([]*domainPayment.InvoicePayment, int, error) {
	var conditions []string
	var args []interface{}

	// Siempre filtrar por invoice_id y excluir soft-deleted
	conditions = append(conditions, "ip.invoice_id = ?")
	args = append(args, invoiceID)
	conditions = append(conditions, "ip.deleted_at IS NULL")

	// Búsqueda por payment_id (fiel al original)
	if search, ok := filters["search"].(string); ok && search != "" {
		conditions = append(conditions, "ip.payment_id LIKE ?")
		args = append(args, "%"+search+"%")
	}

	whereClause := strings.Join(conditions, " AND ")

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM invoice_payments ip
		WHERE %s
	`, whereClause)

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		slog.ErrorContext(ctx, "Failed to count invoice payments",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	// Sorting
	orderClause := "ip.created_at DESC"
	if sort, ok := filters["sort"].(string); ok && sort != "" {
		direction := "DESC"
		if dir, ok := filters["direction"].(string); ok && strings.ToUpper(dir) == "ASC" {
			direction = "ASC"
		}

		switch sort {
		case "amount":
			orderClause = fmt.Sprintf("ip.amount %s", direction)
		case "created_at":
			orderClause = fmt.Sprintf("ip.created_at %s", direction)
		case "payment_processor":
			orderClause = fmt.Sprintf("ip.payment_processor %s", direction)
		}
	}

	// Data query
	offset := (page - 1) * pageSize
	dataQuery := fmt.Sprintf(`
		SELECT
			ip.id, ip.invoice_id, ip.payment_processor, ip.payment_id, ip.amount, ip.notes,
			ip.created_at, ip.updated_at, ip.deleted_at
		FROM invoice_payments ip
		WHERE %s
		ORDER BY %s
		LIMIT ? OFFSET ?
	`, whereClause, orderClause)

	queryArgs := append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, queryArgs...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list invoice payments",
			slog.String("error", err.Error()))
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var payments []*domainPayment.InvoicePayment
	for rows.Next() {
		payment := &domainPayment.InvoicePayment{}
		if err := rows.Scan(
			&payment.ID, &payment.InvoiceID, &payment.PaymentProcessor, &payment.PaymentID, &payment.Amount, &payment.Notes,
			&payment.CreatedAt, &payment.UpdatedAt, &payment.DeletedAt,
		); err != nil {
			slog.ErrorContext(ctx, "Failed to scan invoice payment row",
				slog.String("error", err.Error()))
			return nil, 0, err
		}
		payments = append(payments, payment)
	}

	if err = rows.Err(); err != nil {
		slog.ErrorContext(ctx, "Error iterating invoice payment rows",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	return payments, total, nil
}
