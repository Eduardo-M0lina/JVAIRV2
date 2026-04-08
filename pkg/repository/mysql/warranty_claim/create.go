package warranty_claim

import (
	"context"
	"log/slog"

	domainWC "github.com/angumol/jvairv2/pkg/domain/warranty_claim"
)

func (r *Repository) Create(ctx context.Context, wc *domainWC.WarrantyClaim) error {
	query := `
		INSERT INTO warranty_claims (
			internal_claim_number, warranty_claim_type_id, warranty_claim_status_id, job_id,
			invoice_number, work_done, warranty_part, manufacturer, model_number,
			part_number, replacement_part_number, part_distributor, part_invoice_number,
			old_part_serial_number, new_part_serial_number, esa_number, serial,
			claim_number, approved, parts_credit_received, labor_payment_received,
			notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		wc.InternalClaimNumber,
		wc.WarrantyClaimTypeID,
		wc.WarrantyClaimStatusID,
		wc.JobID,
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
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute insert warranty claim query",
			slog.String("error", err.Error()))
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get last insert ID",
			slog.String("error", err.Error()))
		return err
	}

	wc.ID = id
	return nil
}
