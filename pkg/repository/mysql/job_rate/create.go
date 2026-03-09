package job_rate

import (
	"context"

	"github.com/your-org/jvairv2/pkg/domain/job_rate"
)

func (r *repository) Create(ctx context.Context, rate *job_rate.JobRate) error {
	query := `
		INSERT INTO job_rates (
			job_id, user_id, job_rate_status_id, sale_price, rate_percent, rate_flat,
			tech_parts, company_parts, parts_replaced, deduction, payment, paid, notes,
			created_at, updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		rate.JobID,
		rate.UserID,
		rate.JobRateStatusID,
		rate.SalePrice,
		rate.RatePercent,
		rate.RateFlat,
		rate.TechParts,
		rate.CompanyParts,
		rate.PartsReplaced,
		rate.Deduction,
		rate.Payment,
		rate.Paid,
		rate.Notes,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	rate.ID = id
	return nil
}
