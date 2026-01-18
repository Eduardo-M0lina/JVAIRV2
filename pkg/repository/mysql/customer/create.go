package customer

import (
	"context"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/customer"
)

func (r *Repository) Create(ctx context.Context, c *customer.Customer) error {
	query := `
		INSERT INTO customers (
			name, email, phone, mobile, fax, phone_other, website,
			contact_name, contact_email, contact_phone,
			billing_address_street, billing_address_city, billing_address_state, billing_address_zip,
			workflow_id, notes, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		c.Name,
		c.Email,
		c.Phone,
		c.Mobile,
		c.Fax,
		c.PhoneOther,
		c.Website,
		c.ContactName,
		c.ContactEmail,
		c.ContactPhone,
		c.BillingAddressStreet,
		c.BillingAddressCity,
		c.BillingAddressState,
		c.BillingAddressZip,
		c.WorkflowID,
		c.Notes,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute insert customer query",
			slog.String("error", err.Error()))
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get last insert ID",
			slog.String("error", err.Error()))
		return err
	}

	c.ID = id
	return nil
}
