package customer

import (
	"context"
	"log/slog"
	"strings"

	"github.com/angumol/jvairv2/pkg/domain/customer"
)

func (r *Repository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*customer.Customer, int, error) {
	baseQuery := `
		SELECT
			c.id, c.name, c.email, c.phone, c.mobile, c.fax, c.phone_other, c.website,
			c.contact_name, c.contact_email, c.contact_phone,
			c.billing_address_street, c.billing_address_city, c.billing_address_state, c.billing_address_zip,
			c.workflow_id, c.notes, c.created_at, c.updated_at, c.deleted_at,
			COALESCE(COUNT(p.id), 0) as total_properties
		FROM customers c
		LEFT JOIN properties p ON c.id = p.customer_id AND p.deleted_at IS NULL
		WHERE c.deleted_at IS NULL
	`

	countQuery := "SELECT COUNT(*) FROM customers c WHERE c.deleted_at IS NULL"

	var args []interface{}
	var conditions []string

	if workflowID, ok := filters["workflow_id"].(int64); ok && workflowID > 0 {
		conditions = append(conditions, "c.workflow_id = ?")
		args = append(args, workflowID)
	}

	if billingAddressCity, ok := filters["billing_address_city"].(string); ok && billingAddressCity != "" {
		conditions = append(conditions, "c.billing_address_city = ?")
		args = append(args, billingAddressCity)
	}

	if billingAddressState, ok := filters["billing_address_state"].(string); ok && billingAddressState != "" {
		conditions = append(conditions, "c.billing_address_state = ?")
		args = append(args, billingAddressState)
	}

	if billingAddressZip, ok := filters["billing_address_zip"].(string); ok && billingAddressZip != "" {
		conditions = append(conditions, "c.billing_address_zip = ?")
		args = append(args, billingAddressZip)
	}

	if search, ok := filters["search"].(string); ok && search != "" {
		searchCondition := `(
			c.name LIKE ? OR
			c.email LIKE ? OR
			c.phone LIKE ? OR
			c.mobile LIKE ? OR
			c.contact_name LIKE ? OR
			c.contact_email LIKE ? OR
			c.billing_address_city LIKE ? OR
			c.billing_address_state LIKE ?
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

	baseQuery += " GROUP BY c.id"

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
			&c.TotalProperties,
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
