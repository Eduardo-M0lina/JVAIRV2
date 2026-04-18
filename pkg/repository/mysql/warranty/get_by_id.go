package warranty

import (
	"context"
	"database/sql"
	"log/slog"

	domainWarranty "github.com/angumol/jvairv2/pkg/domain/warranty"
)

func (r *Repository) GetByID(ctx context.Context, id int64) (*domainWarranty.Warranty, error) {
	query := `
		SELECT w.id, w.warranty_number, w.job_id, w.warranty_type_id, w.warranty_status_id,
			w.date_submitted, w.agreement_number, w.audit_done, w.notes,
			w.created_at, w.updated_at, w.deleted_at,
			j.id, j.completion_date,
			p.id, CONCAT(p.street, ', ', p.city, ', ', p.state, ' ', p.zip),
			c.name
		FROM warranties w
		INNER JOIN jobs j ON w.job_id = j.id
		INNER JOIN properties p ON j.property_id = p.id
		INNER JOIN customers c ON p.customer_id = c.id
		WHERE w.id = ? AND w.deleted_at IS NULL
	`

	w := &domainWarranty.Warranty{}
	var dateSubmitted sql.NullTime
	var agreementNumber sql.NullString
	var notes sql.NullString
	var jobID int64
	var completionDate sql.NullTime
	var propertyID int64
	var propertyAddress string
	var customerName string

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
		&jobID, &completionDate,
		&propertyID, &propertyAddress,
		&customerName,
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

	w.Job = &domainWarranty.Job{
		ID: jobID,
		Property: domainWarranty.Property{
			ID:      propertyID,
			Address: propertyAddress,
			Customer: domainWarranty.Customer{
				Name: customerName,
			},
		},
	}

	if completionDate.Valid {
		w.Job.CompletionDate = &completionDate.Time
	}

	return w, nil
}
