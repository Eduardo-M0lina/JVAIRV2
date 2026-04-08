package warranty

import (
	"context"
	"log/slog"

	domainWarranty "github.com/angumol/jvairv2/pkg/domain/warranty"
)

func (r *Repository) Update(ctx context.Context, w *domainWarranty.Warranty) error {
	query := `
		UPDATE warranties
		SET warranty_number = ?, warranty_type_id = ?, warranty_status_id = ?,
			date_submitted = ?, agreement_number = ?, audit_done = ?, notes = ?,
			updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query,
		w.WarrantyNumber,
		w.WarrantyTypeID,
		w.WarrantyStatusID,
		w.DateSubmitted,
		w.AgreementNumber,
		w.AuditDone,
		w.Notes,
		w.ID,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to update warranty",
			slog.String("error", err.Error()),
			slog.Int64("id", w.ID))
		return err
	}

	return nil
}
