package customer

import (
	"context"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/customer"
)

func (r *Repository) Update(ctx context.Context, c *customer.Customer) error {
	query := `
		UPDATE customers SET
			name = ?,
			email = ?,
			phone = ?,
			mobile = ?,
			fax = ?,
			phone_other = ?,
			website = ?,
			contact_name = ?,
			contact_email = ?,
			contact_phone = ?,
			billing_address_street = ?,
			billing_address_city = ?,
			billing_address_state = ?,
			billing_address_zip = ?,
			workflow_id = ?,
			notes = ?,
			updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
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
		c.ID,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to update customer",
			slog.String("error", err.Error()),
			slog.Int64("customer_id", c.ID))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get rows affected",
			slog.String("error", err.Error()))
		return err
	}

	if rowsAffected == 0 {
		slog.WarnContext(ctx, "No rows affected during update",
			slog.Int64("customer_id", c.ID))
	}

	return nil
}
