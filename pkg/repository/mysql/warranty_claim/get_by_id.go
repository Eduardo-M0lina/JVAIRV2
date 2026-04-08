package warranty_claim

import (
	"context"
	"database/sql"
	"log/slog"

	domainWC "github.com/angumol/jvairv2/pkg/domain/warranty_claim"
)

func (r *Repository) GetByID(ctx context.Context, id int64) (*domainWC.WarrantyClaim, error) {
	query := `
		SELECT id, internal_claim_number, warranty_claim_type_id, warranty_claim_status_id, job_id,
			invoice_number, work_done, warranty_part, manufacturer, model_number,
			part_number, replacement_part_number, part_distributor, part_invoice_number,
			old_part_serial_number, new_part_serial_number, esa_number, serial,
			claim_number, approved, parts_credit_received, labor_payment_received,
			notes, created_at, updated_at, deleted_at
		FROM warranty_claims
		WHERE id = ? AND deleted_at IS NULL
	`

	wc := &domainWC.WarrantyClaim{}
	var invoiceNumber, warrantyPart, manufacturer, modelNumber sql.NullString
	var partNumber, replacementPartNumber, partDistributor, partInvoiceNumber sql.NullString
	var oldPartSerialNumber, newPartSerialNumber, esaNumber, serial sql.NullString
	var claimNumber, notes sql.NullString

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

	return wc, nil
}
