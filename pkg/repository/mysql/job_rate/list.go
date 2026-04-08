package job_rate

import (
	"context"

	"github.com/angumol/jvairv2/pkg/domain/job_rate"
)

func (r *repository) List(ctx context.Context, jobID int64, limit, offset int) ([]*job_rate.JobRate, int64, error) {
	query := `
		SELECT id, job_id, user_id, job_rate_status_id, sale_price, rate_percent, rate_flat,
			   tech_parts, company_parts, parts_replaced, deduction, payment, paid, notes,
			   created_at, updated_at, deleted_at
		FROM job_rates
		WHERE job_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, jobID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var rates []*job_rate.JobRate
	for rows.Next() {
		rate := &job_rate.JobRate{}
		err := rows.Scan(
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
		if err != nil {
			return nil, 0, err
		}
		rates = append(rates, rate)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	countQuery := "SELECT COUNT(*) FROM job_rates WHERE job_id = ? AND deleted_at IS NULL"
	var total int64
	err = r.db.QueryRowContext(ctx, countQuery, jobID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return rates, total, nil
}
