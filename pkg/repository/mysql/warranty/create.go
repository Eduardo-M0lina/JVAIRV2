package warranty

import (
	"context"
	"log/slog"

	domainWarranty "github.com/angumol/jvairv2/pkg/domain/warranty"
)

func (r *Repository) Create(ctx context.Context, w *domainWarranty.Warranty) error {
	query := `
		INSERT INTO warranties (warranty_number, job_id, warranty_type_id, warranty_status_id, date_submitted, agreement_number, audit_done, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		w.WarrantyNumber,
		w.JobID,
		w.WarrantyTypeID,
		w.WarrantyStatusID,
		w.DateSubmitted,
		w.AgreementNumber,
		w.AuditDone,
		w.Notes,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute insert warranty query",
			slog.String("error", err.Error()))
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get last insert ID",
			slog.String("error", err.Error()))
		return err
	}

	w.ID = id
	return nil
}
