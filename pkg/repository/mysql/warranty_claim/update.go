package warranty_claim

import (
	"context"
	"log/slog"

	domainWC "github.com/your-org/jvairv2/pkg/domain/warranty_claim"
)

func (r *Repository) Update(ctx context.Context, wc *domainWC.WarrantyClaim) error {
	query := `
		UPDATE warranty_claims
		SET internal_claim_number = ?, warranty_claim_type_id = ?, warranty_claim_status_id = ?,
			invoice_number = ?, work_done = ?, warranty_part = ?, manufacturer = ?, model_number = ?,
			part_number = ?, replacement_part_number = ?, part_distributor = ?, part_invoice_number = ?,
			old_part_serial_number = ?, new_part_serial_number = ?, esa_number = ?, serial = ?,
			claim_number = ?, approved = ?, parts_credit_received = ?, labor_payment_received = ?,
			notes = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query,
		wc.InternalClaimNumber,
		wc.WarrantyClaimTypeID,
		wc.WarrantyClaimStatusID,
		wc.InvoiceNumber,
		wc.WorkDone,
		wc.WarrantyPart,
		wc.Manufacturer,
		wc.ModelNumber,
		wc.PartNumber,
		wc.ReplacementPartNumber,
		wc.PartDistributor,
		wc.PartInvoiceNumber,
		wc.OldPartSerialNumber,
		wc.NewPartSerialNumber,
		wc.EsaNumber,
		wc.Serial,
		wc.ClaimNumber,
		wc.Approved,
		wc.PartsCreditReceived,
		wc.LaborPaymentReceived,
		wc.Notes,
		wc.ID,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to update warranty claim",
			slog.String("error", err.Error()),
			slog.Int64("id", wc.ID))
		return err
	}

	return nil
}
