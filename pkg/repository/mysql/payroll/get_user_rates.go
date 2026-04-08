package payroll

import (
	"context"
	"fmt"
	"strings"

	"github.com/angumol/jvairv2/pkg/domain/payroll"
)

func (r *repository) GetUserRates(ctx context.Context, userID int64, status *string, limit, offset int) ([]*payroll.PayrollRate, int64, error) {
	query := `
		SELECT jr.id, jr.job_id, jr.user_id, jr.job_rate_status_id, jr.sale_price,
			   jr.rate_percent, jr.rate_flat, jr.tech_parts, jr.company_parts,
			   jr.parts_replaced, jr.deduction, jr.payment, jr.paid, jr.notes,
			   jr.created_at, jr.updated_at,
			   j.work_order,
			   CONCAT(p.street, ', ', p.city, ', ', p.state, ' ', p.zip) as property_address,
			   jrs.label as status_label, jrs.class as status_class
		FROM job_rates jr
		LEFT JOIN jobs j ON j.id = jr.job_id
		LEFT JOIN properties p ON p.id = j.property_id
		LEFT JOIN job_rate_statuses jrs ON jrs.id = jr.job_rate_status_id
		WHERE jr.user_id = ? AND jr.deleted_at IS NULL
	`

	var args []interface{}
	args = append(args, userID)

	// Filtrar por status si se especifica
	if status != nil && *status != "" && *status != "all" {
		query += " AND jrs.label = ?"
		// Capitalizar primera letra (reemplazo de strings.Title deprecado)
		statusValue := strings.ToLower(*status)
		if len(statusValue) > 0 {
			statusValue = strings.ToUpper(statusValue[:1]) + statusValue[1:]
		}
		args = append(args, statusValue)
	}

	// Contar total
	countQuery := strings.Replace(query, "SELECT jr.id, jr.job_id, jr.user_id, jr.job_rate_status_id, jr.sale_price, \n\t\t\t   jr.rate_percent, jr.rate_flat, jr.tech_parts, jr.company_parts, \n\t\t\t   jr.parts_replaced, jr.deduction, jr.payment, jr.paid, jr.notes,\n\t\t\t   jr.created_at, jr.updated_at,\n\t\t\t   j.work_order,\n\t\t\t   CONCAT(p.street, ', ', p.city, ', ', p.state, ' ', p.zip) as property_address,\n\t\t\t   jrs.label as status_label, jrs.class as status_class", "SELECT COUNT(*)", 1)

	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting user rates: %w", err)
	}

	// Agregar paginación
	query += " ORDER BY jr.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying user rates: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var rates []*payroll.PayrollRate
	for rows.Next() {
		rate := &payroll.PayrollRate{}
		err := rows.Scan(
			&rate.ID, &rate.JobID, &rate.UserID, &rate.JobRateStatusID,
			&rate.SalePrice, &rate.RatePercent, &rate.RateFlat,
			&rate.TechParts, &rate.CompanyParts, &rate.PartsReplaced,
			&rate.Deduction, &rate.Payment, &rate.Paid, &rate.Notes,
			&rate.CreatedAt, &rate.UpdatedAt,
			&rate.WorkOrder, &rate.PropertyAddr,
			&rate.StatusLabel, &rate.StatusClass,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning rate: %w", err)
		}
		rates = append(rates, rate)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating rates: %w", err)
	}

	return rates, total, nil
}
