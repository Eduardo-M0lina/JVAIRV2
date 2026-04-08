package payroll

import (
	"context"
	"fmt"
	"time"

	"github.com/angumol/jvairv2/pkg/domain/payroll"
)

func (r *repository) GetPaystubData(ctx context.Context, userID int64) (*payroll.PaystubData, error) {
	// Obtener información del usuario
	userQuery := `
		SELECT u.id, u.name, u.email, ro.name as role_name
		FROM users u
		LEFT JOIN assigned_roles ar ON ar.entity_id = u.id AND ar.entity_type = 'App\\Models\\User'
		LEFT JOIN roles ro ON ro.id = ar.role_id
		WHERE u.id = ? AND u.deleted_at IS NULL
		LIMIT 1
	`

	user := &payroll.PayrollUser{}
	err := r.db.QueryRowContext(ctx, userQuery, userID).Scan(
		&user.ID, &user.Name, &user.Email, &user.RoleName,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	// Obtener rates del usuario (solo unpaid y holding para el paystub)
	ratesQuery := `
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
		AND jrs.label IN ('Unpaid', 'Holding')
		ORDER BY jr.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, ratesQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying paystub rates: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var rates []*payroll.PayrollRate
	var totalPayment float64

	user.UnpaidRates = []*payroll.PayrollRate{}
	user.HoldingRates = []*payroll.PayrollRate{}
	user.PaidRates = []*payroll.PayrollRate{}

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
			return nil, fmt.Errorf("error scanning rate: %w", err)
		}
		rates = append(rates, rate)
		totalPayment += rate.Payment

		// Agrupar en el usuario directamente desde la misma query
		if rate.StatusLabel != nil {
			switch *rate.StatusLabel {
			case "Unpaid":
				user.UnpaidRates = append(user.UnpaidRates, rate)
				user.TotalUnpaid += rate.Payment
			case "Holding":
				user.HoldingRates = append(user.HoldingRates, rate)
				user.TotalHolding += rate.Payment
			}
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rates: %w", err)
	}

	return &payroll.PaystubData{
		User:         user,
		Rates:        rates,
		TotalPayment: totalPayment,
		GeneratedAt:  time.Now(),
	}, nil
}
