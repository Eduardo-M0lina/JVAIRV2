package customer

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/customer"
)

func (r *Repository) GetByID(ctx context.Context, id int64) (*customer.Customer, error) {
	query := `
		SELECT
			id, name, email, phone, mobile, fax, phone_other, website,
			contact_name, contact_email, contact_phone,
			billing_address_street, billing_address_city, billing_address_state, billing_address_zip,
			workflow_id, notes, created_at, updated_at, deleted_at
		FROM customers
		WHERE id = ? AND deleted_at IS NULL
	`

	c := &customer.Customer{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
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
		if errors.Is(err, sql.ErrNoRows) {
			slog.WarnContext(ctx, "Customer not found",
				slog.Int64("customer_id", id))
			return nil, errors.New("customer not found")
		}
		slog.ErrorContext(ctx, "Failed to get customer by ID",
			slog.String("error", err.Error()),
			slog.Int64("customer_id", id))
		return nil, err
	}

	return c, nil
}
