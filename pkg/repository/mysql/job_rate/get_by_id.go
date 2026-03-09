package job_rate

import (
	"context"
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/job_rate"
)

func (r *repository) GetByID(ctx context.Context, id int64) (*job_rate.JobRate, error) {
	query := `
		SELECT id, job_id, user_id, job_rate_status_id, sale_price, rate_percent, rate_flat,
			   tech_parts, company_parts, parts_replaced, deduction, payment, paid, notes,
			   created_at, updated_at, deleted_at
		FROM job_rates
		WHERE id = ?
	`

	rate := &job_rate.JobRate{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rate.ID,
		&rate.JobID,
		&rate.UserID,
		&rate.JobRateStatusID,
		&rate.SalePrice,
		&rate.RatePercent,
		&rate.RateFlat,
		&rate.TechParts,
		&rate.CompanyParts,
		&rate.PartsReplaced,
		&rate.Deduction,
		&rate.Payment,
		&rate.Paid,
		&rate.Notes,
		&rate.CreatedAt,
		&rate.UpdatedAt,
		&rate.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, job_rate.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return rate, nil
}
