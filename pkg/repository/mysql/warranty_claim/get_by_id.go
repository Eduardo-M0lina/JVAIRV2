package warranty_claim

import (
	"context"
	"database/sql"
	"log/slog"

	domainWC "github.com/angumol/jvairv2/pkg/domain/warranty_claim"
)

func (r *Repository) GetByID(ctx context.Context, id int64) (*domainWC.WarrantyClaim, error) {
	query := `
		SELECT wc.id, wc.internal_claim_number, wc.warranty_claim_type_id, wc.warranty_claim_status_id, wc.job_id,
			wc.invoice_number, wc.work_done, wc.warranty_part, wc.manufacturer, wc.model_number,
			wc.part_number, wc.replacement_part_number, wc.part_distributor, wc.part_invoice_number,
			wc.old_part_serial_number, wc.new_part_serial_number, wc.esa_number, wc.serial,
			wc.claim_number, wc.approved, wc.parts_credit_received, wc.labor_payment_received,
			wc.notes, wc.created_at, wc.updated_at, wc.deleted_at,
			j.id, j.completion_date,
			p.id, CONCAT(p.street, ', ', p.city, ', ', p.state, ' ', p.zip),
			c.name
		FROM warranty_claims wc
		INNER JOIN jobs j ON wc.job_id = j.id
		INNER JOIN properties p ON j.property_id = p.id
		INNER JOIN customers c ON p.customer_id = c.id
		WHERE wc.id = ? AND wc.deleted_at IS NULL
	`

	wc := &domainWC.WarrantyClaim{}
	var invoiceNumber, warrantyPart, manufacturer, modelNumber sql.NullString
	var partNumber, replacementPartNumber, partDistributor, partInvoiceNumber sql.NullString
	var oldPartSerialNumber, newPartSerialNumber, esaNumber, serial sql.NullString
	var claimNumber, notes sql.NullString
	var jobID int64
	var completionDate sql.NullTime
	var propertyID int64
	var propertyAddress string
	var customerName string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&wc.ID,
		&wc.InternalClaimNumber,
		&wc.WarrantyClaimTypeID,
		&wc.WarrantyClaimStatusID,
		&wc.JobID,
		&invoiceNumber, &wc.WorkDone, &warrantyPart, &manufacturer, &modelNumber,
		&partNumber, &replacementPartNumber, &partDistributor, &partInvoiceNumber,
		&oldPartSerialNumber, &newPartSerialNumber, &esaNumber, &serial,
		&claimNumber, &wc.Approved, &wc.PartsCreditReceived, &wc.LaborPaymentReceived,
		&notes, &wc.CreatedAt, &wc.UpdatedAt, &wc.DeletedAt,
		&jobID, &completionDate,
		&propertyID, &propertyAddress,
		&customerName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domainWC.ErrWarrantyClaimNotFound
		}
		slog.ErrorContext(ctx, "Failed to get warranty claim by ID",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return nil, err
	}

	wc.InvoiceNumber = fromNullString(invoiceNumber)
	wc.WarrantyPart = fromNullString(warrantyPart)
	wc.Manufacturer = fromNullString(manufacturer)
	wc.ModelNumber = fromNullString(modelNumber)
	wc.PartNumber = fromNullString(partNumber)
	wc.ReplacementPartNumber = fromNullString(replacementPartNumber)
	wc.PartDistributor = fromNullString(partDistributor)
	wc.PartInvoiceNumber = fromNullString(partInvoiceNumber)
	wc.OldPartSerialNumber = fromNullString(oldPartSerialNumber)
	wc.NewPartSerialNumber = fromNullString(newPartSerialNumber)
	wc.EsaNumber = fromNullString(esaNumber)
	wc.Serial = fromNullString(serial)
	wc.ClaimNumber = fromNullString(claimNumber)
	wc.Notes = fromNullString(notes)

	wc.Job = &domainWC.Job{
		ID: jobID,
		Property: domainWC.Property{
			ID:      propertyID,
			Address: propertyAddress,
			Customer: domainWC.Customer{
				Name: customerName,
			},
		},
	}

	if completionDate.Valid {
		wc.Job.CompletionDate = &completionDate.Time
	}

	return wc, nil
}
