package customer

import (
	"context"
	"log/slog"
	"strings"

	"github.com/your-org/jvairv2/pkg/domain/customer"
)

func (r *Repository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*customer.Customer, int, error) {
	baseQuery := `
		SELECT
			id, name, email, phone, mobile, fax, phone_other, website,
			contact_name, contact_email, contact_phone,
			billing_address_street, billing_address_city, billing_address_state, billing_address_zip,
			workflow_id, notes, created_at, updated_at, deleted_at
		FROM customers
		WHERE deleted_at IS NULL
	`

	countQuery := "SELECT COUNT(*) FROM customers WHERE deleted_at IS NULL"

	var args []interface{}
	var conditions []string

	if workflowID, ok := filters["workflow_id"].(int64); ok && workflowID > 0 {
		conditions = append(conditions, "workflow_id = ?")
		args = append(args, workflowID)
	}

	if search, ok := filters["search"].(string); ok && search != "" {
		searchCondition := `(
			name LIKE ? OR
			email LIKE ? OR
			phone LIKE ? OR
			mobile LIKE ? OR
			contact_name LIKE ? OR
			contact_email LIKE ? OR
			billing_address_city LIKE ? OR
			billing_address_state LIKE ?
		)`
		conditions = append(conditions, searchCondition)
		searchPattern := "%" + search + "%"
		for i := 0; i < 8; i++ {
			args = append(args, searchPattern)
		}
	}

	if len(conditions) > 0 {
		whereClause := " AND " + strings.Join(conditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
	}

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to count customers",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	baseQuery += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	offset := (page - 1) * pageSize
	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to query customers",
			slog.String("error", err.Error()))
		return nil, 0, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			slog.ErrorContext(ctx, "Failed to close rows", slog.String("error", closeErr.Error()))
		}
	}()

	var customers []*customer.Customer
	for rows.Next() {
		c := &customer.Customer{}
		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Email,
			&c.Phone,
			&c.Mobile,
			&c.Fax,
			&c.PhoneOther,
			&c.Website,
			&c.ContactName,
			&c.ContactEmail,
			&c.ContactPhone,
			&c.BillingAddressStreet,
			&c.BillingAddressCity,
			&c.BillingAddressState,
			&c.BillingAddressZip,
			&c.WorkflowID,
			&c.Notes,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.DeletedAt,
		)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to scan customer row",
				slog.String("error", err.Error()))
			return nil, 0, err
		}
		customers = append(customers, c)
	}

	if err = rows.Err(); err != nil {
		slog.ErrorContext(ctx, "Error iterating customer rows",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	return customers, total, nil
}
