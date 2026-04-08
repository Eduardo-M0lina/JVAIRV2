package warranty

import (
	"context"
	"database/sql"
	"log/slog"

	domainWarranty "github.com/angumol/jvairv2/pkg/domain/warranty"
)

func (r *Repository) GetByID(ctx context.Context, id int64) (*domainWarranty.Warranty, error) {
	query := `
		SELECT id, warranty_number, job_id, warranty_type_id, warranty_status_id,
			date_submitted, agreement_number, audit_done, notes,
			created_at, updated_at, deleted_at
		FROM warranties
		WHERE id = ? AND deleted_at IS NULL
	`

	w := &domainWarranty.Warranty{}
	var dateSubmitted sql.NullTime
	var agreementNumber sql.NullString
	var notes sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&w.ID,
		&w.WarrantyNumber,
		&w.JobID,
		&w.WarrantyTypeID,
		&w.WarrantyStatusID,
		&dateSubmitted,
		&agreementNumber,
		&w.AuditDone,
		&notes,
		&w.CreatedAt,
		&w.UpdatedAt,
		&w.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domainWarranty.ErrWarrantyNotFound
		}
		slog.ErrorContext(ctx, "Failed to get warranty by ID",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return nil, err
	}

	if dateSubmitted.Valid {
		w.DateSubmitted = &dateSubmitted.Time
	}
	if agreementNumber.Valid {
		w.AgreementNumber = &agreementNumber.String
	}
	if notes.Valid {
		w.Notes = &notes.String
	}

	return w, nil
}
