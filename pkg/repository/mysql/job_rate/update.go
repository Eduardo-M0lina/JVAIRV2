package job_rate

import (
	"context"

	"github.com/your-org/jvairv2/pkg/domain/job_rate"
)

func (r *repository) Update(ctx context.Context, rate *job_rate.JobRate) error {
	query := `
		UPDATE job_rates
		SET user_id = ?, job_rate_status_id = ?, sale_price = ?, rate_percent = ?, rate_flat = ?,
			tech_parts = ?, company_parts = ?, parts_replaced = ?, deduction = ?, payment = ?,
			paid = ?, notes = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query,
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
		rate.ID,
	)

	return err
}
